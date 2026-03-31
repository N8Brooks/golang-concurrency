// Package testsuite contains reusable behavioral tests for dining philosophers implementations.
package testsuite

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

type DiningPhilosophers interface {
	Dine(ctx context.Context, philosopher int, think, eat func())
}

type philosopherRun struct {
	thought chan struct{}
	entered chan struct{}
	release chan struct{}
	done    chan struct{}
}

func startPhilosopher(dp DiningPhilosophers, ctx context.Context, philosopher int) philosopherRun {
	run := philosopherRun{
		thought: make(chan struct{}),
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan struct{}),
	}

	go func() {
		dp.Dine(ctx, philosopher, func() {
			close(run.thought)
		}, func() {
			close(run.entered)
			<-run.release
		})
		close(run.done)
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

func requirePending(t *testing.T, done <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-done:
		t.Fatal(message)
	default:
	}
}

func requireSuccess(t *testing.T, done <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-done:
	default:
		t.Fatal(message)
	}
}

func runHungryWave(t *testing.T, dp DiningPhilosophers, ctx context.Context) {
	t.Helper()

	runs := make([]philosopherRun, 5)
	released := make([]bool, 5)
	finished := make([]bool, 5)

	for i := range 5 {
		runs[i] = startPhilosopher(dp, ctx, i)
	}

	synctest.Wait()

	for i := range 5 {
		requireClosed(t, runs[i].thought, fmt.Sprintf("philosopher %d did not think", i))
	}

	completed := 0
	for completed < 5 {
		progress := false

		for i := range 5 {
			if released[i] {
				continue
			}

			select {
			case <-runs[i].entered:
				close(runs[i].release)
				released[i] = true
				progress = true
			default:
			}
		}

		synctest.Wait()

		for i := range 5 {
			if finished[i] {
				continue
			}

			select {
			case <-runs[i].done:
				finished[i] = true
				completed++
				progress = true
			default:
			}
		}

		if completed < 5 && !progress {
			t.Fatal("philosophers deadlocked")
		}
	}
}

func Run(t *testing.T, newImpl func() DiningPhilosophers) {
	t.Helper()

	t.Run("AdjacentPhilosophersConflict", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			dp := newImpl()
			ctx := t.Context()

			first := startPhilosopher(dp, ctx, 0)
			synctest.Wait()

			requireClosed(t, first.thought, "philosopher 0 did not think")
			requireClosed(t, first.entered, "philosopher 0 did not start eating")

			second := startPhilosopher(dp, ctx, 1)
			synctest.Wait()

			requireClosed(t, second.thought, "philosopher 1 did not think")
			requireOpen(t, second.entered, "adjacent philosopher ate while shared fork was unavailable")
			requirePending(t, second.done, "adjacent philosopher completed while waiting for a shared fork")

			close(first.release)
			synctest.Wait()

			requireSuccess(t, first.done, "philosopher 0 did not complete")
			requireClosed(t, second.entered, "philosopher 1 did not start eating after the shared fork was released")

			close(second.release)
			synctest.Wait()

			requireSuccess(t, second.done, "philosopher 1 did not complete")
		})
	})

	t.Run("NonAdjacentPhilosophersCanEatConcurrently", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			dp := newImpl()
			ctx := t.Context()

			first := startPhilosopher(dp, ctx, 0)
			synctest.Wait()

			requireClosed(t, first.entered, "philosopher 0 did not start eating")

			second := startPhilosopher(dp, ctx, 2)
			synctest.Wait()

			requireClosed(t, second.thought, "philosopher 2 did not think")
			requireClosed(t, second.entered, "non-adjacent philosopher did not start eating concurrently")

			close(first.release)
			close(second.release)
			synctest.Wait()

			requireSuccess(t, first.done, "philosopher 0 did not complete")
			requireSuccess(t, second.done, "philosopher 2 did not complete")
		})
	})

	t.Run("AvoidsDeadlockWhenAllHungry", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runHungryWave(t, newImpl(), t.Context())
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			dp := newImpl()
			ctx := t.Context()

			for range 5 {
				runHungryWave(t, dp, ctx)
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			dp := newImpl()

			first := startPhilosopher(dp, t.Context(), 0)
			synctest.Wait()

			requireClosed(t, first.entered, "philosopher 0 did not start eating")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			second := startPhilosopher(dp, ctx, 1)
			synctest.Wait()

			requireClosed(t, second.thought, "philosopher 1 did not think")
			requireOpen(t, second.entered, "philosopher 1 ate despite the shared fork being unavailable")
			requirePending(t, second.done, "philosopher 1 completed before cancellation")

			cancel()
			synctest.Wait()

			select {
			case <-second.done:
			default:
				t.Fatal("philosopher 1 did not exit after cancellation")
			}

			requireOpen(t, second.entered, "philosopher 1 ate despite being canceled while waiting")

			close(first.release)
			synctest.Wait()

			requireSuccess(t, first.done, "philosopher 0 did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() DiningPhilosophers) {
	b.Helper()
	b.ReportAllocs()
	empty := func() {}

	b.Run("Contended", func(b *testing.B) {
		ctx := b.Context()
		dp := newImpl()
		var next atomic.Uint32

		b.RunParallel(func(pb *testing.PB) {
			philosopher := int(next.Add(1)-1) % 5
			for pb.Next() {
				ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
				dp.Dine(ctx, philosopher, empty, empty)
				cancel()
			}
		})
	})
}
