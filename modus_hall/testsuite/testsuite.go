// Package testsuite contains reusable behavioral tests for Modus Hall implementations.
package testsuite

import (
	"context"
	"errors"
	"sync"
	"testing"
	"testing/synctest"
)

type ModusHall interface {
	Heathen(ctx context.Context, cross func()) error
	Prude(ctx context.Context, cross func()) error
}

type worker struct {
	entered chan struct{}
	release chan struct{}
	done    chan error
}

func startHeathen(mh ModusHall, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- mh.Heathen(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()

	return w
}

func startPrude(mh ModusHall, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- mh.Prude(ctx, func() {
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

func Run(t *testing.T, newImpl func() ModusHall) {
	t.Helper()

	t.Run("HeathensShare", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mh := newImpl()
			ctx := t.Context()

			h1 := startHeathen(mh, ctx)
			synctest.Wait()
			requireClosed(t, h1.entered, "first heathen did not enter")

			h2 := startHeathen(mh, ctx)
			synctest.Wait()
			requireClosed(t, h2.entered, "second heathen did not enter alongside first heathen")

			close(h1.release)
			close(h2.release)
			synctest.Wait()

			requireSuccess(t, h1.done, "first heathen did not complete")
			requireSuccess(t, h2.done, "second heathen did not complete")
		})
	})

	t.Run("PrudesShare", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mh := newImpl()
			ctx := t.Context()

			p1 := startPrude(mh, ctx)
			synctest.Wait()
			requireClosed(t, p1.entered, "first prude did not enter")

			p2 := startPrude(mh, ctx)
			synctest.Wait()
			requireClosed(t, p2.entered, "second prude did not enter alongside first prude")

			close(p1.release)
			close(p2.release)
			synctest.Wait()

			requireSuccess(t, p1.done, "first prude did not complete")
			requireSuccess(t, p2.done, "second prude did not complete")
		})
	})

	t.Run("OppositionWaitsForField", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mh := newImpl()
			ctx := t.Context()

			h := startHeathen(mh, ctx)
			synctest.Wait()
			requireClosed(t, h.entered, "heathen did not enter")

			p := startPrude(mh, ctx)
			synctest.Wait()
			requireOpen(t, p.entered, "prude entered while a heathen was active")

			close(h.release)
			synctest.Wait()

			requireSuccess(t, h.done, "heathen did not complete")
			requireClosed(t, p.entered, "prude did not enter after the field cleared")

			close(p.release)
			synctest.Wait()
			requireSuccess(t, p.done, "prude did not complete")
		})
	})

	t.Run("MajorityTransitionBlocksIncumbents", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mh := newImpl()
			ctx := t.Context()

			h1 := startHeathen(mh, ctx)
			h2 := startHeathen(mh, ctx)
			synctest.Wait()
			requireClosed(t, h1.entered, "first heathen did not enter")
			requireClosed(t, h2.entered, "second heathen did not enter")

			p1 := startPrude(mh, ctx)
			p2 := startPrude(mh, ctx)
			synctest.Wait()
			requireOpen(t, p1.entered, "first prude entered before the field cleared")
			requireOpen(t, p2.entered, "second prude entered before the field cleared")

			lateHeathen := startHeathen(mh, ctx)
			synctest.Wait()
			requireClosed(t, lateHeathen.entered, "late heathen should still enter before a prude majority exists")

			p3 := startPrude(mh, ctx)
			p4 := startPrude(mh, ctx)
			synctest.Wait()
			requireOpen(t, p3.entered, "third prude entered before the field cleared")
			requireOpen(t, p4.entered, "fourth prude entered before the field cleared")

			blockedHeathen := startHeathen(mh, ctx)
			synctest.Wait()
			requireOpen(t, blockedHeathen.entered, "new heathen entered after prudes gained majority")

			close(h1.release)
			close(h2.release)
			close(lateHeathen.release)
			synctest.Wait()

			requireSuccess(t, h1.done, "first heathen did not complete")
			requireSuccess(t, h2.done, "second heathen did not complete")
			requireSuccess(t, lateHeathen.done, "late heathen did not complete")
			requireClosed(t, p1.entered, "first prude did not enter after the heathens drained")
			requireClosed(t, p2.entered, "second prude did not enter after the heathens drained")
			requireClosed(t, p3.entered, "third prude did not enter after the heathens drained")
			requireClosed(t, p4.entered, "fourth prude did not enter after the heathens drained")
			requireOpen(t, blockedHeathen.entered, "blocked heathen entered before the prudes finished")

			close(p1.release)
			close(p2.release)
			close(p3.release)
			close(p4.release)
			synctest.Wait()

			requireSuccess(t, p1.done, "first prude did not complete")
			requireSuccess(t, p2.done, "second prude did not complete")
			requireSuccess(t, p3.done, "third prude did not complete")
			requireSuccess(t, p4.done, "fourth prude did not complete")
			requireClosed(t, blockedHeathen.entered, "blocked heathen did not enter after the prudes finished")

			close(blockedHeathen.release)
			synctest.Wait()
			requireSuccess(t, blockedHeathen.done, "blocked heathen did not complete")
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			setup func(mh ModusHall, ctx context.Context) worker
			start func(mh ModusHall, ctx context.Context) worker
		}{
			{name: "Heathen", setup: startPrude, start: startHeathen},
			{name: "Prude", setup: startHeathen, start: startPrude},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					mh := newImpl()
					ctx := t.Context()

					active := tc.setup(mh, ctx)
					synctest.Wait()
					requireClosed(t, active.entered, "active worker did not enter")

					waitCtx, cancel := context.WithCancel(ctx)
					defer cancel()

					waiter := tc.start(mh, waitCtx)
					synctest.Wait()
					requireOpen(t, waiter.entered, tc.name+" entered despite the field being occupied")

					cancel()
					synctest.Wait()
					requireCanceled(t, waiter.done, tc.name+" did not exit after cancellation")

					close(active.release)
					synctest.Wait()
					requireSuccess(t, active.done, "active worker did not complete")
				})
			})
		}
	})

	t.Run("ReusableHandoffs", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mh := newImpl()
			ctx := t.Context()

			h := startHeathen(mh, ctx)
			synctest.Wait()
			requireClosed(t, h.entered, "heathen did not enter")

			p1 := startPrude(mh, ctx)
			p2 := startPrude(mh, ctx)
			synctest.Wait()
			requireOpen(t, p1.entered, "first prude entered before heathen left")
			requireOpen(t, p2.entered, "second prude entered before heathen left")

			close(h.release)
			synctest.Wait()
			requireSuccess(t, h.done, "heathen did not complete")
			requireClosed(t, p1.entered, "first prude did not enter after heathen left")
			requireClosed(t, p2.entered, "second prude did not enter after heathen left")

			nextHeathen := startHeathen(mh, ctx)
			synctest.Wait()
			requireOpen(t, nextHeathen.entered, "next heathen entered while prudes were active")

			close(p1.release)
			close(p2.release)
			synctest.Wait()
			requireSuccess(t, p1.done, "first prude did not complete")
			requireSuccess(t, p2.done, "second prude did not complete")
			requireClosed(t, nextHeathen.entered, "next heathen did not enter after prudes left")

			close(nextHeathen.release)
			synctest.Wait()
			requireSuccess(t, nextHeathen.done, "next heathen did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() ModusHall) {
	b.Helper()
	b.ReportAllocs()

	b.Run("HeathenBatch", func(b *testing.B) {
		for b.Loop() {
			mh := newImpl()
			ctx := b.Context()

			var wg sync.WaitGroup
			wg.Add(4)
			for range 4 {
				go func() {
					defer wg.Done()
					if err := mh.Heathen(ctx, func() {}); err != nil {
						b.Errorf("heathen returned %v", err)
					}
				}()
			}
			wg.Wait()
		}
	})

	b.Run("MajorityHandoff", func(b *testing.B) {
		for b.Loop() {
			mh := newImpl()
			ctx := b.Context()

			var wg sync.WaitGroup
			wg.Add(5)

			for range 2 {
				go func() {
					defer wg.Done()
					if err := mh.Heathen(ctx, func() {}); err != nil {
						b.Errorf("heathen returned %v", err)
					}
				}()
			}

			for range 3 {
				go func() {
					defer wg.Done()
					if err := mh.Prude(ctx, func() {}); err != nil {
						b.Errorf("prude returned %v", err)
					}
				}()
			}

			wg.Wait()
		}
	})
}
