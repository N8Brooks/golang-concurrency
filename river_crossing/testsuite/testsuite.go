// Package testsuite contains reusable behavioral tests for river crossing implementations.
package testsuite

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type RiverCrossing interface {
	Hacker(ctx context.Context, board, rowBoat func()) error
	Serf(ctx context.Context, board, rowBoat func()) error
}

type passengerKind int

const (
	hacker passengerKind = iota
	serf
)

func (k passengerKind) String() string {
	switch k {
	case hacker:
		return "Hacker"
	case serf:
		return "Serf"
	default:
		return "Unknown"
	}
}

type participant struct {
	kind passengerKind
}

func runParticipant(r RiverCrossing, ctx context.Context, p passengerKind, board, rowBoat func()) error {
	if p == hacker {
		return r.Hacker(ctx, board, rowBoat)
	}
	return r.Serf(ctx, board, rowBoat)
}

func requireClosed(t *testing.T, ch <-chan error, message string) error {
	t.Helper()

	select {
	case err := <-ch:
		return err
	default:
		t.Fatal(message)
		return nil
	}
}

func requireNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

func runOneBoat(t *testing.T, r RiverCrossing, ctx context.Context, crew []passengerKind) {
	t.Helper()

	var hackersBoarded atomic.Int32
	var serfsBoarded atomic.Int32
	var boarded atomic.Int32
	var rowed atomic.Int32

	done := make([]chan error, len(crew))
	for i, kind := range crew {
		done[i] = make(chan error, 1)
		go func(ch chan error, kind passengerKind) {
			ch <- runParticipant(r, ctx, kind, func() {
				switch kind {
				case hacker:
					hackersBoarded.Add(1)
				case serf:
					serfsBoarded.Add(1)
				}
				if state := boarded.Add(1); state > 4 {
					t.Errorf("board called %d times, want 4", state)
				}
				if rowed.Load() != 0 {
					t.Error("board ran after rowBoat")
				}
			}, func() {
				if boarded.Load() != 4 {
					t.Errorf("rowBoat ran after %d board calls, want 4", boarded.Load())
				}
				if rowed.Add(1) != 1 {
					t.Error("rowBoat called more than once for one boatload")
				}
			})
		}(done[i], kind)
		synctest.Wait()
	}

	synctest.Wait()

	for i, ch := range done {
		requireNoError(t, requireClosed(t, ch, crew[i].String()+" did not complete"), crew[i].String()+" returned error")
	}

	if boarded.Load() != 4 {
		t.Fatalf("board called %d times, want 4", boarded.Load())
	}
	if rowed.Load() != 1 {
		t.Fatalf("rowBoat called %d times, want 1", rowed.Load())
	}

	switch {
	case hackersBoarded.Load() == 4 && serfsBoarded.Load() == 0:
	case hackersBoarded.Load() == 0 && serfsBoarded.Load() == 4:
	case hackersBoarded.Load() == 2 && serfsBoarded.Load() == 2:
	default:
		t.Fatalf("unsafe boatload: hackers=%d serfs=%d", hackersBoarded.Load(), serfsBoarded.Load())
	}
}

func Run(t *testing.T, newImpl func() RiverCrossing) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			crew []passengerKind
		}{
			{name: "FourHackers", crew: []passengerKind{hacker, hacker, hacker, hacker}},
			{name: "FourSerfs", crew: []passengerKind{serf, serf, serf, serf}},
			{name: "MixedAlternating", crew: []passengerKind{hacker, serf, hacker, serf}},
			{name: "MixedSerfFirst", crew: []passengerKind{serf, hacker, serf, hacker}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneBoat(t, newImpl(), t.Context(), tc.crew)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			r := newImpl()
			ctx := t.Context()
			for _, crew := range [][]passengerKind{
				{hacker, hacker, hacker, hacker},
				{serf, hacker, serf, hacker},
				{serf, serf, serf, serf},
				{hacker, serf, hacker, serf},
			} {
				runOneBoat(t, r, ctx, crew)
			}
		})
	})

	t.Run("UnsafeThreeOneDoesNotBoard", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			r := newImpl()
			var hackersBoarded atomic.Int32
			var serfsBoarded atomic.Int32
			var rowed atomic.Int32

			done := make([]chan error, 5)
			crew := []passengerKind{hacker, hacker, hacker, serf, serf}
			for i, kind := range crew {
				done[i] = make(chan error, 1)
				go func(ch chan error, kind passengerKind) {
					ch <- runParticipant(r, ctx, kind, func() {
						if kind == hacker {
							hackersBoarded.Add(1)
						} else {
							serfsBoarded.Add(1)
						}
					}, func() {
						rowed.Add(1)
					})
				}(done[i], kind)
				synctest.Wait()
				if i == 3 {
					if hackersBoarded.Load() != 0 || serfsBoarded.Load() != 0 || rowed.Load() != 0 {
						t.Fatalf("unsafe 3-1 group boarded: hackers=%d serfs=%d rowed=%d", hackersBoarded.Load(), serfsBoarded.Load(), rowed.Load())
					}
				}
			}

			synctest.Wait()

			completed := 0
			for _, ch := range done {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("participant returned %v", err)
					}
					completed++
				default:
				}
			}
			if completed != 4 {
				t.Fatalf("completed %d participants, want 4 after forming a 2-2 crew", completed)
			}
			if hackersBoarded.Load() != 2 || serfsBoarded.Load() != 2 {
				t.Fatalf("got hackers=%d serfs=%d, want hackers=2 serfs=2", hackersBoarded.Load(), serfsBoarded.Load())
			}
			if rowed.Load() != 1 {
				t.Fatalf("rowBoat called %d times, want 1", rowed.Load())
			}

			cancel()
			synctest.Wait()

			remaining := 0
			for _, ch := range done {
				select {
				case err := <-ch:
					if errors.Is(err, context.Canceled) {
						remaining++
					}
				default:
				}
			}
			if remaining != 1 {
				t.Fatalf("got %d canceled leftover participants, want 1", remaining)
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			kind passengerKind
		}{
			{name: "Hacker", kind: hacker},
			{name: "Serf", kind: serf},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()

					r := newImpl()
					boarded := atomic.Bool{}
					rowed := atomic.Bool{}
					done := make(chan error, 1)

					go func() {
						done <- runParticipant(r, ctx, tc.kind, func() {
							boarded.Store(true)
						}, func() {
							rowed.Store(true)
						})
					}()

					synctest.Wait()
					cancel()
					synctest.Wait()

					err := requireClosed(t, done, tc.name+" did not exit after cancellation")
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("%s returned %v, want context.Canceled", tc.name, err)
					}
					if boarded.Load() {
						t.Fatalf("%s boarded despite never finding a crew", tc.name)
					}
					if rowed.Load() {
						t.Fatalf("%s rowed despite never finding a crew", tc.name)
					}
				})
			})
		}
	})

	t.Run("CancelRemovesWaiter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			r := newImpl()

			canceledCtx, cancel := context.WithCancel(t.Context())
			defer cancel()

			canceledDone := make(chan error, 1)
			go func() {
				canceledDone <- r.Hacker(canceledCtx, func() {}, func() {
					t.Error("canceled hacker rowed")
				})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			err := requireClosed(t, canceledDone, "canceled hacker did not exit")
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled hacker returned %v, want context.Canceled", err)
			}

			done := make([]chan error, 4)
			crew := []passengerKind{hacker, serf, serf, hacker}
			for i, kind := range crew {
				done[i] = make(chan error, 1)
				go func(ch chan error, kind passengerKind) {
					ch <- runParticipant(r, t.Context(), kind, func() {}, func() {})
				}(done[i], kind)
			}

			synctest.Wait()

			for i, ch := range done {
				requireNoError(t, requireClosed(t, ch, crew[i].String()+" did not complete"), crew[i].String()+" returned error")
			}
		})
	})

	t.Run("BoardingDoesNotInterleaveAcrossBoats", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			r := newImpl()
			ctx := t.Context()

			firstBoatBoarded := make(chan struct{}, 4)
			var secondBoatBoarded atomic.Int32
			releaseFirstBoat := make(chan struct{})

			firstDone := make([]chan error, 4)
			for i := range 4 {
				firstDone[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- r.Hacker(ctx, func() {
						firstBoatBoarded <- struct{}{}
						<-releaseFirstBoat
					}, func() {})
				}(firstDone[i])
			}

			for range 4 {
				<-firstBoatBoarded
			}

			secondDone := make([]chan error, 4)
			for i := range 4 {
				secondDone[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- r.Serf(ctx, func() {
						secondBoatBoarded.Add(1)
					}, func() {})
				}(secondDone[i])
			}

			if secondBoatBoarded.Load() != 0 {
				t.Fatalf("second boat boarded %d participants before the first boat completed", secondBoatBoarded.Load())
			}

			close(releaseFirstBoat)

			for _, ch := range firstDone {
				if err := <-ch; err != nil {
					t.Fatalf("first-boat passenger returned error: %v", err)
				}
			}
			for _, ch := range secondDone {
				if err := <-ch; err != nil {
					t.Fatalf("second-boat passenger returned error: %v", err)
				}
			}

			if secondBoatBoarded.Load() != 4 {
				t.Fatalf("second boat boarded %d participants after the first boat completed, want 4", secondBoatBoarded.Load())
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() RiverCrossing) {
	b.Helper()
	b.ReportAllocs()

	for _, tc := range []struct {
		name string
		crew []passengerKind
	}{
		{name: "FourHackers", crew: []passengerKind{hacker, hacker, hacker, hacker}},
		{name: "Mixed", crew: []passengerKind{hacker, serf, hacker, serf}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				ctx := b.Context()
				r := newImpl()
				done := make([]chan error, len(tc.crew))
				for i, kind := range tc.crew {
					done[i] = make(chan error, 1)
					go func(ch chan error, kind passengerKind) {
						ch <- runParticipant(r, ctx, kind, func() {}, func() {})
					}(done[i], kind)
				}
				for i, ch := range done {
					if err := <-ch; err != nil {
						b.Fatalf("%s passenger %d returned %v", tc.crew[i].String(), i, err)
					}
				}
			}
		})
	}
}
