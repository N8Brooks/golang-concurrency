// Package testsuite contains reusable behavioral tests for multiplex implementations.
package testsuite

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
)

type Multiplex interface {
	Run(ctx context.Context, criticalSection func())
}

type worker struct {
	entered chan struct{}
	release chan struct{}
	done    chan struct{}
}

func startWorker(m Multiplex, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan struct{}),
	}

	go func() {
		m.Run(ctx, func() {
			close(w.entered)
			<-w.release
		})
		close(w.done)
	}()

	return w
}

func requireClosed(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
	default:
		t.Fatal(message)
	}
}

func requireOpen(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
		t.Fatal(message)
	default:
	}
}

func runWave(t *testing.T, m Multiplex, ctx context.Context, limit int) {
	t.Helper()

	workers := make([]worker, limit+1)

	for i := range limit {
		workers[i] = startWorker(m, ctx)
		synctest.Wait()
		requireClosed(t, workers[i].entered, fmt.Sprintf("worker %d did not enter", i))
	}

	workers[limit] = startWorker(m, ctx)
	synctest.Wait()

	requireOpen(t, workers[limit].entered, "extra worker entered despite full capacity")
	requireOpen(t, workers[limit].done, "extra worker exited before a slot became available")

	close(workers[0].release)
	synctest.Wait()

	requireClosed(t, workers[0].done, "released worker did not complete")
	requireClosed(t, workers[limit].entered, "extra worker did not enter after a slot was released")

	for i := 1; i <= limit; i++ {
		close(workers[i].release)
	}

	synctest.Wait()

	for i, w := range workers {
		requireClosed(t, w.done, fmt.Sprintf("worker %d did not complete", i))
	}
}

func Run(t *testing.T, newImpl func(limit int) Multiplex) {
	t.Helper()

	t.Run("EnforcesLimit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runWave(t, newImpl(2), t.Context(), 2)
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m := newImpl(3)
			ctx := t.Context()
			for range 5 {
				runWave(t, m, ctx, 3)
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m := newImpl(1)

			first := startWorker(m, t.Context())
			synctest.Wait()
			requireClosed(t, first.entered, "first worker did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			second := startWorker(m, ctx)
			synctest.Wait()

			requireOpen(t, second.entered, "second worker entered despite full capacity")
			requireOpen(t, second.done, "second worker exited before cancellation")

			cancel()
			synctest.Wait()

			requireClosed(t, second.done, "second worker did not exit after cancellation")
			requireOpen(t, second.entered, "second worker entered despite being canceled while waiting")

			close(first.release)
			synctest.Wait()

			requireClosed(t, first.done, "first worker did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(limit int) Multiplex) {
	b.Helper()
	b.ReportAllocs()

	for _, limit := range []int{1, 8} {
		b.Run(fmt.Sprintf("Limit%d", limit), func(b *testing.B) {
			m := newImpl(limit)
			ctx := b.Context()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					m.Run(ctx, func() {})
				}
			})
		})
	}
}
