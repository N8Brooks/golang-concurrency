// Package testsuite contains reusable behavioral tests for dining savages implementations.
package testsuite

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const defaultCapacity = 3

type DiningSavages interface {
	Savage(ctx context.Context, getServingFromPot, eat func())
	Cook(ctx context.Context, putServingsInPot func(servings int))
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

func Run(t *testing.T, newImpl func(capacity int) DiningSavages) {
	t.Helper()

	t.Run("OneSavageMultipleRounds", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const totalMeals = 8

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			ds := newImpl(defaultCapacity)

			var potServings atomic.Int64
			var potOps atomic.Int32
			var meals atomic.Int64
			var fills atomic.Int64
			var hasServing atomic.Bool
			done := make(chan struct{})
			stop := sync.OnceFunc(func() {
				cancel()
				close(done)
			})

			go ds.Cook(ctx, func(servings int) {
				if servings != defaultCapacity {
					t.Errorf("cook filled %d servings, want %d", servings, defaultCapacity)
				}
				if potOps.Add(1) != 1 {
					t.Error("cook refilled the pot concurrently with another pot operation")
				}
				if old := potServings.Swap(int64(servings)); old != 0 {
					t.Errorf("cook refilled the pot while %d servings remained", old)
				}
				fills.Add(1)
				if potOps.Add(-1) != 0 {
					t.Error("pot operation count did not return to zero after refill")
				}
			})

			go ds.Savage(ctx, func() {
				if potOps.Add(1) != 1 {
					t.Error("savage took a serving concurrently with another pot operation")
				}
				if remaining := potServings.Add(-1); remaining < 0 {
					t.Error("savage took a serving from an empty pot")
				}
				if hasServing.Swap(true) {
					t.Error("savage took a serving while already holding one")
				}
				if potOps.Add(-1) != 0 {
					t.Error("pot operation count did not return to zero after getServingFromPot")
				}
			}, func() {
				if !hasServing.Swap(false) {
					t.Error("savage ate without first taking a serving")
				}
				if meals.Add(1) == totalMeals {
					stop()
				}
			})

			<-done
			synctest.Wait()

			if meals.Load() != totalMeals {
				t.Fatalf("savage ate %d meals, want %d", meals.Load(), totalMeals)
			}

			wantFills := int64((totalMeals + defaultCapacity - 1) / defaultCapacity)
			if fills.Load() != wantFills {
				t.Fatalf("cook refilled %d times, want %d", fills.Load(), wantFills)
			}
		})
	})

	t.Run("GetServingFromPotIsExclusive", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ds := newImpl(2)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			go ds.Cook(ctx, func(int) {})

			firstCtx, cancelFirst := context.WithCancel(ctx)
			defer cancelFirst()

			firstGot := make(chan struct{})
			releaseFirstGet := make(chan struct{})
			firstAte := make(chan struct{})
			secondGot := make(chan struct{})
			secondAte := make(chan struct{})

			go ds.Savage(firstCtx, func() {
				close(firstGot)
				<-releaseFirstGet
			}, func() {
				close(firstAte)
				cancelFirst()
			})

			<-firstGot

			go ds.Savage(ctx, func() {
				close(secondGot)
			}, func() {
				close(secondAte)
				cancel()
			})

			requireOpen(t, secondGot, "second savage entered getServingFromPot while the first held the pot")

			close(releaseFirstGet)

			<-firstAte
			<-secondGot
			<-secondAte
		})
	})

	t.Run("EatRunsOutsidePotCriticalSection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ds := newImpl(2)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			go ds.Cook(ctx, func(int) {})

			firstEating := make(chan struct{})
			releaseFirstEat := make(chan struct{})
			secondGot := make(chan struct{})
			secondAte := make(chan struct{})

			go ds.Savage(ctx, func() {}, func() {
				close(firstEating)
				<-releaseFirstEat
			})

			<-firstEating

			go ds.Savage(ctx, func() {
				close(secondGot)
			}, func() {
				close(secondAte)
				cancel()
			})

			<-secondGot

			close(releaseFirstEat)
			<-secondAte
		})
	})

	t.Run("CancelWhileWaitingForCook", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			ds := newImpl(defaultCapacity)

			var gotServing atomic.Int64
			var ate atomic.Int64
			exited := make(chan struct{})

			go func() {
				ds.Savage(ctx, func() {
					gotServing.Add(1)
				}, func() {
					ate.Add(1)
				})
				close(exited)
			}()

			synctest.Wait()

			cancel()
			synctest.Wait()

			requireClosed(t, exited, "savage did not exit after cancellation while waiting for the cook")

			if gotServing.Load() != 0 {
				t.Fatalf("savage got %d servings while waiting for the cook, want 0", gotServing.Load())
			}
			if ate.Load() != 0 {
				t.Fatalf("savage ate %d meals while waiting for the cook, want 0", ate.Load())
			}
		})
	})

	t.Run("CancelIdleCook", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			ds := newImpl(defaultCapacity)
			exited := make(chan struct{})

			go func() {
				ds.Cook(ctx, func(int) {})
				close(exited)
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			requireClosed(t, exited, "cook did not exit after cancellation while idle")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) DiningSavages) {
	b.Helper()
	b.ReportAllocs()

	for _, capacity := range []int{1, 8} {
		b.Run(fmt.Sprintf("Capacity%d", capacity), func(b *testing.B) {
			for b.Loop() {
				ctx, cancel := context.WithCancel(b.Context())
				ds := newImpl(capacity)
				done := make(chan struct{})

				go ds.Cook(ctx, func(int) {})
				go ds.Savage(ctx, func() {}, func() {
					close(done)
					cancel()
				})

				<-done
			}
		})
	}
}
