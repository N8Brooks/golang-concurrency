// Package testsuite contains reusable behavioral tests for sushi bar implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
)

const defaultCapacity = 5

type SushiBar interface {
	Customer(ctx context.Context, dine func()) error
}

type customerRun struct {
	seated  chan struct{}
	release chan struct{}
	done    chan error
}

func startCustomer(bar SushiBar, ctx context.Context) customerRun {
	run := customerRun{
		seated:  make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		run.done <- bar.Customer(ctx, func() {
			close(run.seated)
			<-run.release
		})
	}()

	return run
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

func requireSuccess(t *testing.T, done <-chan error, message string) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("%s: %v", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func requireCanceled(t *testing.T, done <-chan error, message string) {
	t.Helper()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: got %v, want context.Canceled", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func runOneShot(t *testing.T, bar SushiBar, ctx context.Context, capacity int) {
	t.Helper()

	customers := make([]customerRun, capacity)
	for i := range capacity {
		customers[i] = startCustomer(bar, ctx)
		synctest.Wait()
		requireClosed(t, customers[i].seated, "customer did not sit while seats were available")
	}

	waiting := startCustomer(bar, ctx)
	synctest.Wait()
	requireOpen(t, waiting.seated, "waiting customer sat while the current party was still full")

	close(customers[0].release)
	synctest.Wait()
	requireSuccess(t, customers[0].done, "first customer did not complete")
	requireOpen(t, waiting.seated, "waiting customer sat after only one diner left")

	for i := 1; i < capacity; i++ {
		close(customers[i].release)
	}
	synctest.Wait()
	for i := 1; i < capacity; i++ {
		requireSuccess(t, customers[i].done, "customer did not complete")
	}

	requireClosed(t, waiting.seated, "waiting customer did not sit after the whole party left")
	close(waiting.release)
	synctest.Wait()
	requireSuccess(t, waiting.done, "waiting customer did not complete")
}

func Run(t *testing.T, newImpl func(capacity int) SushiBar) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runOneShot(t, newImpl(defaultCapacity), t.Context(), defaultCapacity)
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			bar := newImpl(defaultCapacity)
			ctx := t.Context()
			for range 3 {
				runOneShot(t, bar, ctx, defaultCapacity)
			}
		})
	})

	t.Run("AdmitsAtMostCapacityFromWaitingParty", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			bar := newImpl(defaultCapacity)
			ctx := t.Context()

			current := make([]customerRun, defaultCapacity)
			for i := range defaultCapacity {
				current[i] = startCustomer(bar, ctx)
			}
			synctest.Wait()
			for i := range defaultCapacity {
				requireClosed(t, current[i].seated, "current diner did not sit")
			}

			waiting := make([]customerRun, defaultCapacity+1)
			for i := range waiting {
				waiting[i] = startCustomer(bar, ctx)
			}
			synctest.Wait()
			for i := range waiting {
				requireOpen(t, waiting[i].seated, "waiting diner sat before the current party left")
			}

			for i := range current {
				close(current[i].release)
			}
			synctest.Wait()
			for i := range current {
				requireSuccess(t, current[i].done, "current diner did not complete")
			}

			seated := 0
			var extra customerRun
			extraIndex := -1
			for i, customer := range waiting {
				select {
				case <-customer.seated:
					seated++
				default:
					extra = customer
					extraIndex = i
				}
			}

			if seated != defaultCapacity {
				t.Fatalf("got %d seated waiting customers, want %d", seated, defaultCapacity)
			}
			if extraIndex == -1 {
				t.Fatal("expected one waiting customer to remain blocked")
			}

			for i, customer := range waiting {
				if i == extraIndex {
					continue
				}
				close(customer.release)
			}
			synctest.Wait()
			for i, customer := range waiting {
				if i == extraIndex {
					continue
				}
				requireSuccess(t, customer.done, "released waiting diner did not complete")
			}

			requireClosed(t, extra.seated, "extra waiting diner did not sit after the admitted party left")
			close(extra.release)
			synctest.Wait()
			requireSuccess(t, extra.done, "extra waiting diner did not complete")
		})
	})

	t.Run("CancelWaitingCustomer", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			bar := newImpl(defaultCapacity)

			current := make([]customerRun, defaultCapacity)
			for i := range defaultCapacity {
				current[i] = startCustomer(bar, t.Context())
			}
			synctest.Wait()
			for i := range defaultCapacity {
				requireClosed(t, current[i].seated, "current diner did not sit")
			}

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waiting := startCustomer(bar, ctx)
			synctest.Wait()
			requireOpen(t, waiting.seated, "waiting diner sat while the bar was full")

			cancel()
			synctest.Wait()
			requireCanceled(t, waiting.done, "waiting diner did not exit after cancellation")

			for i := range current {
				close(current[i].release)
			}
			synctest.Wait()
			for i := range current {
				requireSuccess(t, current[i].done, "current diner did not complete")
			}

			replacement := startCustomer(bar, t.Context())
			synctest.Wait()
			requireClosed(t, replacement.seated, "replacement diner did not sit after canceled waiter left")
			close(replacement.release)
			synctest.Wait()
			requireSuccess(t, replacement.done, "replacement diner did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) SushiBar) {
	b.Helper()
	b.ReportAllocs()

	for _, capacity := range []int{1, 5} {
		b.Run(fmt.Sprintf("Capacity%d", capacity), func(b *testing.B) {
			bar := newImpl(capacity)
			ctx := b.Context()
			for b.Loop() {
				customer := startCustomer(bar, ctx)
				<-customer.seated
				close(customer.release)
				if err := <-customer.done; err != nil {
					b.Fatalf("customer returned %v", err)
				}
			}
		})
	}
}
