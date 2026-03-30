// Package testsuite contains reusable behavioral tests for mutex implementations.
package testsuite

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock()
}

func runOneShot(t *testing.T, m Mutex, ctx context.Context) {
	t.Helper()

	var inside atomic.Int32
	holderEntered := make(chan struct{})
	holderRelease := make(chan struct{})
	holderDone := make(chan struct{})
	contenderEntered := make(chan struct{})
	contenderDone := make(chan struct{})

	go func() {
		if err := m.Lock(ctx); err != nil {
			t.Errorf("holder lock failed: %v", err)
			close(holderDone)
			return
		}
		if inside.Add(1) != 1 {
			t.Error("holder entered critical section concurrently")
		}
		close(holderEntered)
		<-holderRelease
		if inside.Add(-1) != 0 {
			t.Error("holder left invalid critical section count")
		}
		m.Unlock()
		close(holderDone)
	}()

	synctest.Wait()

	select {
	case <-holderEntered:
	default:
		t.Fatal("holder did not enter the critical section")
	}

	go func() {
		if err := m.Lock(ctx); err != nil {
			t.Errorf("contender lock failed: %v", err)
			close(contenderDone)
			return
		}
		if inside.Add(1) != 1 {
			t.Error("contender entered critical section concurrently")
		}
		close(contenderEntered)
		if inside.Add(-1) != 0 {
			t.Error("contender left invalid critical section count")
		}
		m.Unlock()
		close(contenderDone)
	}()

	synctest.Wait()

	select {
	case <-contenderEntered:
		t.Fatal("contender entered while holder still owned the mutex")
	default:
	}

	close(holderRelease)
	synctest.Wait()

	select {
	case <-holderDone:
	default:
		t.Fatal("holder did not complete")
	}

	select {
	case <-contenderEntered:
	default:
		t.Fatal("contender did not enter after the mutex was released")
	}

	select {
	case <-contenderDone:
	default:
		t.Fatal("contender did not complete")
	}
}

func Run(t *testing.T, newImpl func() Mutex) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runOneShot(t, newImpl(), t.Context())
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m := newImpl()
			ctx := t.Context()
			for range 5 {
				runOneShot(t, m, ctx)
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m := newImpl()

			holderEntered := make(chan struct{})
			holderRelease := make(chan struct{})
			holderDone := make(chan struct{})

			go func() {
				if err := m.Lock(t.Context()); err != nil {
					t.Errorf("holder lock failed: %v", err)
					close(holderDone)
					return
				}
				close(holderEntered)
				<-holderRelease
				m.Unlock()
				close(holderDone)
			}()

			synctest.Wait()

			select {
			case <-holderEntered:
			default:
				t.Fatal("holder did not acquire the mutex")
			}

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waiterDone := make(chan error, 1)

			go func() {
				waiterDone <- m.Lock(ctx)
			}()

			synctest.Wait()

			cancel()
			synctest.Wait()

			select {
			case err := <-waiterDone:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("waiter returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("waiter did not exit after cancellation")
			}

			close(holderRelease)
			synctest.Wait()

			select {
			case <-holderDone:
			default:
				t.Fatal("holder did not complete")
			}

			if err := m.Lock(t.Context()); err != nil {
				t.Fatalf("mutex was not reusable after cancellation: %v", err)
			}
			m.Unlock()
		})
	})
}

func Benchmark(b *testing.B, newImpl func() Mutex) {
	b.Helper()
	b.ReportAllocs()

	b.Run("Uncontended", func(b *testing.B) {
		ctx := b.Context()
		m := newImpl()

		for b.Loop() {
			if err := m.Lock(ctx); err != nil {
				b.Fatalf("lock failed: %v", err)
			}
			m.Unlock()
		}
	})

	b.Run("Contended", func(b *testing.B) {
		ctx := b.Context()

		for b.Loop() {
			m := newImpl()
			holderEntered := make(chan struct{})
			holderRelease := make(chan struct{})
			holderDone := make(chan struct{})
			waiterDone := make(chan struct{})

			go func() {
				if err := m.Lock(ctx); err != nil {
					close(holderDone)
					return
				}
				close(holderEntered)
				<-holderRelease
				m.Unlock()
				close(holderDone)
			}()

			<-holderEntered

			go func() {
				if err := m.Lock(ctx); err == nil {
					m.Unlock()
				}
				close(waiterDone)
			}()

			close(holderRelease)
			<-waiterDone
			<-holderDone
		}
	})
}
