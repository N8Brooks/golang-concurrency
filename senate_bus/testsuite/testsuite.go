// Package testsuite contains reusable behavioral tests for Senate bus implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type SenateBus interface {
	Rider(ctx context.Context, boardBus func()) error
	Bus(ctx context.Context, depart func()) error
}

type riderRun struct {
	boarded chan struct{}
	release chan struct{}
	done    chan error
}

type busRun struct {
	departed chan struct{}
	done     chan error
}

func startRider(bus SenateBus, ctx context.Context) riderRun {
	run := riderRun{
		boarded: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		run.done <- bus.Rider(ctx, func() {
			close(run.boarded)
			<-run.release
		})
	}()

	return run
}

func startBus(bus SenateBus, ctx context.Context) busRun {
	run := busRun{
		departed: make(chan struct{}),
		done:     make(chan error, 1),
	}

	go func() {
		run.done <- bus.Bus(ctx, func() {
			close(run.departed)
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

func requirePending(t *testing.T, done <-chan error, message string) {
	t.Helper()

	select {
	case err := <-done:
		t.Fatalf("%s: returned %v", message, err)
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

func countBoarded(riders ...riderRun) (count int, indexes []int) {
	for i, rider := range riders {
		select {
		case <-rider.boarded:
			count++
			indexes = append(indexes, i)
		default:
		}
	}
	return count, indexes
}

func Run(t *testing.T, newImpl func(capacity int) SenateBus) {
	t.Helper()

	t.Run("EmptyBusDepartsImmediately", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			bus := startBus(newImpl(3), t.Context())
			synctest.Wait()

			requireClosed(t, bus.departed, "bus did not depart immediately when no riders were waiting")
			requireSuccess(t, bus.done, "bus did not complete")
		})
	})

	t.Run("WaitingRidersBoardThenBusDeparts", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sb := newImpl(3)
			ctx := t.Context()

			r1 := startRider(sb, ctx)
			r2 := startRider(sb, ctx)
			synctest.Wait()

			requireOpen(t, r1.boarded, "first rider boarded before a bus arrived")
			requireOpen(t, r2.boarded, "second rider boarded before a bus arrived")

			bus := startBus(sb, ctx)
			synctest.Wait()

			requireClosed(t, r1.boarded, "first rider did not board when the bus arrived")
			requireClosed(t, r2.boarded, "second rider did not board when the bus arrived")
			requireOpen(t, bus.departed, "bus departed before riders finished boarding")
			requirePending(t, bus.done, "bus returned before riders finished boarding")

			close(r1.release)
			close(r2.release)
			synctest.Wait()

			requireSuccess(t, r1.done, "first rider did not complete")
			requireSuccess(t, r2.done, "second rider did not complete")
			requireClosed(t, bus.departed, "bus did not depart after riders boarded")
			requireSuccess(t, bus.done, "bus did not complete")
		})
	})

	t.Run("CapacityLimit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sb := newImpl(3)
			ctx := t.Context()

			riders := []riderRun{
				startRider(sb, ctx),
				startRider(sb, ctx),
				startRider(sb, ctx),
				startRider(sb, ctx),
				startRider(sb, ctx),
			}
			released := make([]bool, len(riders))
			completed := make([]bool, len(riders))
			synctest.Wait()

			bus1 := startBus(sb, ctx)
			synctest.Wait()

			count, boarded := countBoarded(riders...)
			if count != 3 {
				t.Fatalf("got %d riders on the first bus, want 3", count)
			}
			requireOpen(t, bus1.departed, "bus departed before boarded riders finished")

			for _, i := range boarded {
				close(riders[i].release)
				released[i] = true
			}
			synctest.Wait()

			for _, i := range boarded {
				requireSuccess(t, riders[i].done, fmt.Sprintf("rider %d did not complete on first bus", i))
				completed[i] = true
			}
			requireClosed(t, bus1.departed, "first bus did not depart")
			requireSuccess(t, bus1.done, "first bus did not complete")

			count, boarded = countBoarded(riders...)
			if count != 3 {
				t.Fatalf("got %d total boarded riders after first departure, want 3", count)
			}
			bus2 := startBus(sb, ctx)
			synctest.Wait()

			count, boarded = countBoarded(riders...)
			if count != 5 {
				t.Fatalf("got %d total boarded riders after second bus, want 5", count)
			}

			for _, i := range boarded {
				if released[i] {
					continue
				}
				close(riders[i].release)
				released[i] = true
			}
			synctest.Wait()

			for i, rider := range riders {
				if completed[i] {
					continue
				}
				requireSuccess(t, rider.done, fmt.Sprintf("rider %d did not complete", i))
			}
			requireClosed(t, bus2.departed, "second bus did not depart")
			requireSuccess(t, bus2.done, "second bus did not complete")
		})
	})

	t.Run("LateArrivalsWaitForNextBus", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sb := newImpl(3)
			ctx := t.Context()

			r1 := startRider(sb, ctx)
			r2 := startRider(sb, ctx)
			synctest.Wait()

			bus1 := startBus(sb, ctx)
			synctest.Wait()

			requireClosed(t, r1.boarded, "first rider did not board first bus")
			requireClosed(t, r2.boarded, "second rider did not board first bus")
			requireOpen(t, bus1.departed, "first bus departed before riders finished boarding")

			late := startRider(sb, ctx)
			synctest.Wait()
			requireOpen(t, late.boarded, "late rider boarded the current bus")

			close(r1.release)
			close(r2.release)
			synctest.Wait()

			requireSuccess(t, r1.done, "first rider did not complete")
			requireSuccess(t, r2.done, "second rider did not complete")
			requireClosed(t, bus1.departed, "first bus did not depart")
			requireSuccess(t, bus1.done, "first bus did not complete")
			requireOpen(t, late.boarded, "late rider boarded before next bus arrived")

			bus2 := startBus(sb, ctx)
			synctest.Wait()

			requireClosed(t, late.boarded, "late rider did not board the next bus")
			requireOpen(t, bus2.departed, "second bus departed before late rider finished boarding")

			close(late.release)
			synctest.Wait()

			requireSuccess(t, late.done, "late rider did not complete")
			requireClosed(t, bus2.departed, "second bus did not depart")
			requireSuccess(t, bus2.done, "second bus did not complete")
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sb := newImpl(2)
			ctx := t.Context()

			for range 5 {
				r1 := startRider(sb, ctx)
				r2 := startRider(sb, ctx)
				synctest.Wait()

				bus := startBus(sb, ctx)
				synctest.Wait()

				close(r1.release)
				close(r2.release)
				synctest.Wait()

				requireSuccess(t, r1.done, "first rider did not complete")
				requireSuccess(t, r2.done, "second rider did not complete")
				requireClosed(t, bus.departed, "bus did not depart")
				requireSuccess(t, bus.done, "bus did not complete")
			}
		})
	})

	t.Run("CancelWaitingRider", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sb := newImpl(2)

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			canceled := startRider(sb, ctx)
			synctest.Wait()

			requireOpen(t, canceled.boarded, "rider boarded before a bus arrived")
			cancel()
			synctest.Wait()
			requireCanceled(t, canceled.done, "rider did not exit after cancellation")

			remaining := startRider(sb, t.Context())
			synctest.Wait()

			bus := startBus(sb, t.Context())
			synctest.Wait()

			requireClosed(t, remaining.boarded, "remaining rider did not board")
			close(remaining.release)
			synctest.Wait()

			requireSuccess(t, remaining.done, "remaining rider did not complete")
			requireClosed(t, bus.departed, "bus did not depart")
			requireSuccess(t, bus.done, "bus did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) SenateBus) {
	b.Helper()
	b.ReportAllocs()

	for _, capacity := range []int{1, 8} {
		b.Run(fmt.Sprintf("Capacity%d", capacity), func(b *testing.B) {
			for b.Loop() {
				sb := newImpl(capacity)
				ctx := b.Context()
				var wg sync.WaitGroup
				var boarded atomic.Int64
				wg.Add(capacity)

				for range capacity {
					go func() {
						defer wg.Done()
						if err := sb.Rider(ctx, func() {}); err != nil {
							b.Errorf("rider returned %v", err)
							return
						}
						boarded.Add(1)
					}()
				}

				for boarded.Load() < int64(capacity) {
					if err := sb.Bus(ctx, func() {}); err != nil {
						b.Errorf("bus returned %v", err)
						break
					}
				}

				wg.Wait()
			}
		})
	}
}
