// Package testsuite contains reusable behavioral tests for child-care
// implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
)

type ChildCare interface {
	Child(ctx context.Context, child func()) error
	Adult(ctx context.Context, adult func()) error
}

type worker struct {
	entered chan struct{}
	release chan struct{}
	done    chan error
}

func startChild(c ChildCare, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- c.Child(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()

	return w
}

func startAdult(c ChildCare, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- c.Adult(ctx, func() {
			close(w.entered)
			<-w.release
		})
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

func requireSuccess(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("%s: %v", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func requireCanceled(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: got %v, want context.Canceled", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func requirePending(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		t.Fatalf("%s: got early completion %v", message, err)
	default:
	}
}

func runCapacityWave(t *testing.T, c ChildCare, ctx context.Context) {
	t.Helper()

	adult := startAdult(c, ctx)
	synctest.Wait()
	requireClosed(t, adult.entered, "adult did not enter")

	children := make([]worker, 4)
	for i := range 3 {
		children[i] = startChild(c, ctx)
		synctest.Wait()
		requireClosed(t, children[i].entered, fmt.Sprintf("child %d did not enter", i))
	}

	children[3] = startChild(c, ctx)
	synctest.Wait()
	requireOpen(t, children[3].entered, "fourth child entered despite only one adult being inside")
	requirePending(t, children[3].done, "fourth child completed before capacity was available")

	close(children[0].release)
	synctest.Wait()
	requireSuccess(t, children[0].done, "released child did not complete")
	requireClosed(t, children[3].entered, "fourth child did not enter after a child left")

	for i := 1; i < len(children); i++ {
		close(children[i].release)
	}
	synctest.Wait()

	for i := 1; i < len(children); i++ {
		requireSuccess(t, children[i].done, fmt.Sprintf("child %d did not complete", i))
	}

	close(adult.release)
	synctest.Wait()
	requireSuccess(t, adult.done, "adult did not complete")
}

func Run(t *testing.T, newImpl func() ChildCare) {
	t.Helper()

	t.Run("ChildWaitsWithoutAdult", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := newImpl()

			child := startChild(c, t.Context())
			synctest.Wait()
			requireOpen(t, child.entered, "child entered without an adult present")
			requirePending(t, child.done, "child completed without an adult present")

			adult := startAdult(c, t.Context())
			synctest.Wait()
			requireClosed(t, adult.entered, "adult did not enter")
			requireClosed(t, child.entered, "child did not enter after an adult arrived")

			close(child.release)
			synctest.Wait()
			requireSuccess(t, child.done, "child did not complete")

			close(adult.release)
			synctest.Wait()
			requireSuccess(t, adult.done, "adult did not complete")
		})
	})

	t.Run("OneAdultSupportsThreeChildren", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runCapacityWave(t, newImpl(), t.Context())
		})
	})

	t.Run("AdultArrivalAddsCapacity", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := newImpl()
			ctx := t.Context()

			firstAdult := startAdult(c, ctx)
			synctest.Wait()
			requireClosed(t, firstAdult.entered, "first adult did not enter")

			children := make([]worker, 4)
			for i := range 3 {
				children[i] = startChild(c, ctx)
				synctest.Wait()
				requireClosed(t, children[i].entered, fmt.Sprintf("child %d did not enter", i))
			}

			children[3] = startChild(c, ctx)
			synctest.Wait()
			requireOpen(t, children[3].entered, "fourth child entered before another adult arrived")

			secondAdult := startAdult(c, ctx)
			synctest.Wait()
			requireClosed(t, secondAdult.entered, "second adult did not enter")
			requireClosed(t, children[3].entered, "fourth child did not enter after the second adult arrived")

			for i := range children {
				close(children[i].release)
			}
			synctest.Wait()

			for i, child := range children {
				requireSuccess(t, child.done, fmt.Sprintf("child %d did not complete", i))
			}

			close(firstAdult.release)
			close(secondAdult.release)
			synctest.Wait()
			requireSuccess(t, firstAdult.done, "first adult did not complete")
			requireSuccess(t, secondAdult.done, "second adult did not complete")
		})
	})

	t.Run("AdultLeaveIsAtomicAcrossThreePermits", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := newImpl()
			ctx := t.Context()

			adults := []worker{
				startAdult(c, ctx),
				startAdult(c, ctx),
			}
			synctest.Wait()
			for i, adult := range adults {
				requireClosed(t, adult.entered, fmt.Sprintf("adult %d did not enter", i))
			}

			children := make([]worker, 3)
			for i := range children {
				children[i] = startChild(c, ctx)
				synctest.Wait()
				requireClosed(t, children[i].entered, fmt.Sprintf("child %d did not enter", i))
			}

			close(adults[0].release)
			close(adults[1].release)
			synctest.Wait()

			completed := 0
			waitingAdult := -1
			for i, adult := range adults {
				select {
				case err := <-adult.done:
					if err != nil {
						t.Fatalf("adult %d returned %v", i, err)
					}
					completed++
				default:
					waitingAdult = i
				}
			}

			if completed != 1 {
				t.Fatalf("got %d adults that left, want exactly 1", completed)
			}
			if waitingAdult < 0 {
				t.Fatal("expected one adult to remain inside")
			}

			close(children[0].release)
			synctest.Wait()
			requireSuccess(t, children[0].done, "first child did not complete")
			requirePending(t, adults[waitingAdult].done, "adult left after only one child exited")

			close(children[1].release)
			synctest.Wait()
			requireSuccess(t, children[1].done, "second child did not complete")
			requirePending(t, adults[waitingAdult].done, "adult left after only two children exited")

			close(children[2].release)
			synctest.Wait()
			requireSuccess(t, children[2].done, "third child did not complete")
			requireSuccess(t, adults[waitingAdult].done, "remaining adult did not complete after all children exited")
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := newImpl()
			ctx := t.Context()
			for range 3 {
				runCapacityWave(t, c, ctx)
			}
		})
	})

	t.Run("CancelWhileChildWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := newImpl()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waitingChild := startChild(c, ctx)
			synctest.Wait()
			requireOpen(t, waitingChild.entered, "child entered before an adult arrived")
			requirePending(t, waitingChild.done, "child completed before cancellation")

			cancel()
			synctest.Wait()
			requireCanceled(t, waitingChild.done, "child did not exit after cancellation")
			requireOpen(t, waitingChild.entered, "child entered despite being canceled while waiting")

			adult := startAdult(c, t.Context())
			synctest.Wait()
			requireClosed(t, adult.entered, "adult did not enter after a canceled child")

			child := startChild(c, t.Context())
			synctest.Wait()
			requireClosed(t, child.entered, "new child did not enter after a canceled waiter was removed")

			close(child.release)
			synctest.Wait()
			requireSuccess(t, child.done, "new child did not complete")

			close(adult.release)
			synctest.Wait()
			requireSuccess(t, adult.done, "adult did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() ChildCare) {
	b.Helper()
	b.ReportAllocs()

	for _, adults := range []int{1, 2} {
		b.Run(fmt.Sprintf("Adults%d", adults), func(b *testing.B) {
			c := newImpl()
			ctx := context.Background()

			releases := make([]chan struct{}, adults)
			done := make([]chan error, adults)

			for i := range adults {
				ready := make(chan struct{})
				releases[i] = make(chan struct{})
				done[i] = make(chan error, 1)

				go func(i int) {
					done[i] <- c.Adult(ctx, func() {
						close(ready)
						<-releases[i]
					})
				}(i)

				<-ready
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if err := c.Child(ctx, func() {}); err != nil {
						b.Fatalf("child returned %v", err)
					}
				}
			})
			b.StopTimer()

			for _, release := range releases {
				close(release)
			}
			for i, adultDone := range done {
				if err := <-adultDone; err != nil {
					b.Fatalf("adult %d returned %v", i, err)
				}
			}
		})
	}
}
