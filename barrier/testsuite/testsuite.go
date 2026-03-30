// Package testsuite contains reusable behavioral tests for barrier implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const defaultParties = 4

type Barrier interface {
	Wait(ctx context.Context, phase1, phase2 func()) error
}

func runOneRound(t *testing.T, b Barrier, ctx context.Context, order []int) {
	t.Helper()

	parties := len(order)
	phase1Mask := uint64(1<<parties) - 1
	phase2Mask := phase1Mask << parties

	var state atomic.Uint64
	done := make([]chan error, parties)

	for _, id := range order {
		done[id] = make(chan error, 1)
		go func(id int) {
			done[id] <- b.Wait(ctx, func() {
				if state.Or(1<<id)&phase2Mask != 0 {
					t.Errorf("phase1 callback %d ran after a phase2 callback", id)
				}
			}, func() {
				if state.Or(1<<(parties+id))&phase1Mask != phase1Mask {
					t.Errorf("phase2 callback %d ran before all phase1 callbacks completed", id)
				}
			})
		}(id)
		synctest.Wait()
	}

	synctest.Wait()

	for _, id := range order {
		select {
		case err := <-done[id]:
			if err != nil {
				t.Fatalf("participant %d returned %v", id, err)
			}
		default:
			t.Fatalf("participant %d did not complete", id)
		}
	}

	if state.Load()&phase1Mask != phase1Mask {
		t.Fatal("not all phase1 callbacks ran")
	}
	if state.Load()&phase2Mask != phase2Mask {
		t.Fatal("not all phase2 callbacks ran")
	}
}

func Run(t *testing.T, newImpl func(parties int) Barrier) {
	t.Helper()

	ascending := []int{0, 1, 2, 3}
	descending := []int{3, 2, 1, 0}

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			order []int
		}{
			{name: "Ascending", order: ascending},
			{name: "Descending", order: descending},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneRound(t, newImpl(defaultParties), t.Context(), tc.order)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			order []int
		}{
			{name: "Ascending", order: ascending},
			{name: "Descending", order: descending},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					b := newImpl(defaultParties)
					ctx := t.Context()
					for range 5 {
						runOneRound(t, b, ctx, tc.order)
					}
				})
			})
		}
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			b := newImpl(defaultParties)

			var state atomic.Uint32
			done := make(chan error, 1)

			go func() {
				done <- b.Wait(ctx, func() {
					state.Or(1)
				}, func() {
					state.Or(2)
				})
			}()

			synctest.Wait()

			cancel()
			synctest.Wait()

			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("wait returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("participant did not exit after cancellation")
			}

			if state.Load()&1 == 0 {
				t.Fatal("phase1 did not run")
			}
			if state.Load()&2 != 0 {
				t.Fatal("phase2 ran despite the round never completing")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func(parties int) Barrier) {
	b.Helper()
	b.ReportAllocs()

	for _, parties := range []int{2, 8} {
		b.Run(fmt.Sprintf("Parties%d", parties), func(b *testing.B) {
			for b.Loop() {
				ctx := b.Context()
				barrier := newImpl(parties)
				ready := make(chan struct{}, parties-1)
				done := make([]chan error, parties-1)

				for i := range parties - 1 {
					done[i] = make(chan error, 1)
					go func(ch chan error) {
						ch <- barrier.Wait(ctx, func() {
							ready <- struct{}{}
						}, func() {})
					}(done[i])
				}

				for range parties - 1 {
					<-ready
				}

				if err := barrier.Wait(ctx, func() {}, func() {}); err != nil {
					b.Fatalf("final participant returned %v", err)
				}

				for i, ch := range done {
					if err := <-ch; err != nil {
						b.Fatalf("participant %d returned %v", i, err)
					}
				}
			}
		})
	}
}
