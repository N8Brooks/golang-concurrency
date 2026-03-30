// Package testsuite contains reusable behavioral tests for no-starve unisex
// bathroom implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type NoStarveUnisexBathroom interface {
	Male(ctx context.Context, bathroom func()) error
	Female(ctx context.Context, bathroom func()) error
}

type worker struct {
	entered chan struct{}
	release chan struct{}
	done    chan error
}

func startMale(b NoStarveUnisexBathroom, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- b.Male(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()

	return w
}

func startFemale(b NoStarveUnisexBathroom, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- b.Female(ctx, func() {
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

func runCapacityWave(t *testing.T, b NoStarveUnisexBathroom, ctx context.Context, male bool) {
	t.Helper()

	workers := make([]worker, 4)
	for i := range 3 {
		if male {
			workers[i] = startMale(b, ctx)
		} else {
			workers[i] = startFemale(b, ctx)
		}
		synctest.Wait()
		requireClosed(t, workers[i].entered, fmt.Sprintf("worker %d did not enter", i))
	}

	if male {
		workers[3] = startMale(b, ctx)
	} else {
		workers[3] = startFemale(b, ctx)
	}
	synctest.Wait()

	requireOpen(t, workers[3].entered, "fourth same-sex worker entered despite bathroom being full")
	requirePending(t, workers[3].done, "fourth same-sex worker exited before space was available")

	close(workers[0].release)
	synctest.Wait()

	requireSuccess(t, workers[0].done, "released worker did not complete")
	requireClosed(t, workers[3].entered, "fourth same-sex worker did not enter after a slot opened")

	for i := 1; i < 4; i++ {
		close(workers[i].release)
	}
	synctest.Wait()

	for i, w := range workers[1:] {
		requireSuccess(t, w.done, fmt.Sprintf("worker %d did not complete", i+1))
	}
}

func Run(t *testing.T, newImpl func() NoStarveUnisexBathroom) {
	t.Helper()

	t.Run("MenShareUpToLimit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runCapacityWave(t, newImpl(), t.Context(), true)
		})
	})

	t.Run("WomenShareUpToLimit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runCapacityWave(t, newImpl(), t.Context(), false)
		})
	})

	t.Run("WomenWaitForMen", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()
			ctx := t.Context()

			male1 := startMale(b, ctx)
			male2 := startMale(b, ctx)
			synctest.Wait()
			requireClosed(t, male1.entered, "first male did not enter")
			requireClosed(t, male2.entered, "second male did not enter")

			female := startFemale(b, ctx)
			synctest.Wait()
			requireOpen(t, female.entered, "female entered while men were inside")

			close(male1.release)
			close(male2.release)
			synctest.Wait()

			requireSuccess(t, male1.done, "first male did not complete")
			requireSuccess(t, male2.done, "second male did not complete")
			requireClosed(t, female.entered, "female did not enter after men left")

			close(female.release)
			synctest.Wait()
			requireSuccess(t, female.done, "female did not complete")
		})
	})

	t.Run("MenWaitForWomen", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()
			ctx := t.Context()

			female1 := startFemale(b, ctx)
			female2 := startFemale(b, ctx)
			synctest.Wait()
			requireClosed(t, female1.entered, "first female did not enter")
			requireClosed(t, female2.entered, "second female did not enter")

			male := startMale(b, ctx)
			synctest.Wait()
			requireOpen(t, male.entered, "male entered while women were inside")

			close(female1.release)
			close(female2.release)
			synctest.Wait()

			requireSuccess(t, female1.done, "first female did not complete")
			requireSuccess(t, female2.done, "second female did not complete")
			requireClosed(t, male.entered, "male did not enter after women left")

			close(male.release)
			synctest.Wait()
			requireSuccess(t, male.done, "male did not complete")
		})
	})

	t.Run("QueuedOppositeSexBlocksLaterSameSex", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()
			ctx := t.Context()

			activeMale := startMale(b, ctx)
			synctest.Wait()
			requireClosed(t, activeMale.entered, "active male did not enter")

			waitingFemale := startFemale(b, ctx)
			synctest.Wait()
			requireOpen(t, waitingFemale.entered, "female entered while male was inside")

			lateMale := startMale(b, ctx)
			synctest.Wait()
			requireOpen(t, lateMale.entered, "late male bypassed queued female")

			close(activeMale.release)
			synctest.Wait()

			requireSuccess(t, activeMale.done, "active male did not complete")
			requireClosed(t, waitingFemale.entered, "queued female did not enter after men left")
			requireOpen(t, lateMale.entered, "late male entered before queued female finished")

			close(waitingFemale.release)
			synctest.Wait()

			requireSuccess(t, waitingFemale.done, "queued female did not complete")
			requireClosed(t, lateMale.entered, "late male did not enter after queued female finished")

			close(lateMale.release)
			synctest.Wait()
			requireSuccess(t, lateMale.done, "late male did not complete")
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()
			ctx := t.Context()
			runCapacityWave(t, b, ctx, true)
			runCapacityWave(t, b, ctx, false)
		})
	})

	t.Run("CancelWhileWaitingForOppositeSex", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()

			active := startFemale(b, t.Context())
			synctest.Wait()
			requireClosed(t, active.entered, "female did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waitingMale := startMale(b, ctx)
			synctest.Wait()
			requireOpen(t, waitingMale.entered, "male entered while female was inside")

			cancel()
			synctest.Wait()
			requireCanceled(t, waitingMale.done, "male did not exit after cancellation")

			close(active.release)
			synctest.Wait()
			requireSuccess(t, active.done, "female did not complete")
		})
	})

	t.Run("CancelWhileWaitingForCapacity", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()
			ctx := t.Context()

			workers := make([]worker, 3)
			for i := range workers {
				workers[i] = startMale(b, ctx)
				synctest.Wait()
				requireClosed(t, workers[i].entered, fmt.Sprintf("male %d did not enter", i))
			}

			canceledCtx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waiting := startMale(b, canceledCtx)
			synctest.Wait()
			requireOpen(t, waiting.entered, "fourth male entered despite bathroom being full")

			cancel()
			synctest.Wait()
			requireCanceled(t, waiting.done, "waiting male did not exit after cancellation")

			replacement := startMale(b, t.Context())
			synctest.Wait()
			requireOpen(t, replacement.entered, "replacement male entered before a slot opened")

			close(workers[0].release)
			synctest.Wait()

			requireSuccess(t, workers[0].done, "released male did not complete")
			requireClosed(t, replacement.entered, "replacement male did not enter after a slot opened")

			close(replacement.release)
			close(workers[1].release)
			close(workers[2].release)
			synctest.Wait()

			requireSuccess(t, replacement.done, "replacement male did not complete")
			requireSuccess(t, workers[1].done, "male 1 did not complete")
			requireSuccess(t, workers[2].done, "male 2 did not complete")
		})
	})

	t.Run("CanceledQueuedOppositeSexDoesNotBlockFutureSameSex", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newImpl()

			activeMale := startMale(b, t.Context())
			synctest.Wait()
			requireClosed(t, activeMale.entered, "active male did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waitingFemale := startFemale(b, ctx)
			synctest.Wait()
			requireOpen(t, waitingFemale.entered, "female entered while male was inside")

			cancel()
			synctest.Wait()
			requireCanceled(t, waitingFemale.done, "queued female did not exit after cancellation")

			lateMale := startMale(b, t.Context())
			synctest.Wait()
			requireClosed(t, lateMale.entered, "late male did not enter after canceled female was removed")

			close(activeMale.release)
			close(lateMale.release)
			synctest.Wait()

			requireSuccess(t, activeMale.done, "active male did not complete")
			requireSuccess(t, lateMale.done, "late male did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() NoStarveUnisexBathroom) {
	b.Helper()
	b.ReportAllocs()

	b.Run("MaleOnly", func(b *testing.B) {
		bathroom := newImpl()
		ctx := b.Context()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := bathroom.Male(ctx, func() {}); err != nil {
					b.Fatalf("male returned %v", err)
				}
			}
		})
	})

	b.Run("Mixed", func(b *testing.B) {
		bathroom := newImpl()
		ctx := b.Context()
		var counter atomic.Uint64

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if counter.Add(1)%2 == 0 {
					if err := bathroom.Male(ctx, func() {}); err != nil {
						b.Fatalf("male returned %v", err)
					}
					continue
				}
				if err := bathroom.Female(ctx, func() {}); err != nil {
					b.Fatalf("female returned %v", err)
				}
			}
		})
	})
}
