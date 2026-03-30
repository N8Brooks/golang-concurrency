package rendezvous_test

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/N8Brooks/golang-concurrency/rendezvous"
)

const (
	a1Bit uint32 = 1 << iota
	b1Bit
	a2Bit
	b2Bit
)

func TestRendezvous(t *testing.T) {
	for _, tc := range []struct {
		name   string
		aFirst bool
	}{
		{name: "a-first", aFirst: true},
		{name: "b-first", aFirst: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				r := rendezvous.NewRendezvous()
				ctx := t.Context()

				var state atomic.Uint32

				doneA := make(chan struct{})
				doneB := make(chan struct{})

				a1 := func() {
					if state.Or(a1Bit)&(a2Bit|b2Bit) != 0 {
						t.Error("a1 ran after a second phase callback")
					}
				}
				a2 := func() {
					if state.Or(a2Bit)&(a1Bit|b1Bit) != (a1Bit | b1Bit) {
						t.Error("a2 ran before both first phase callbacks completed")
					}
				}
				b1 := func() {
					if state.Or(b1Bit)&(a2Bit|b2Bit) != 0 {
						t.Error("b1 ran after a second phase callback")
					}
				}
				b2 := func() {
					if state.Or(b2Bit)&(a1Bit|b1Bit) != (a1Bit | b1Bit) {
						t.Error("b2 ran before both first phase callbacks completed")
					}
				}

				runA := func() {
					r.A(ctx, a1, a2)
					close(doneA)
				}
				runB := func() {
					r.B(ctx, b1, b2)
					close(doneB)
				}

				if tc.aFirst {
					go runA()
					synctest.Wait()
					go runB()
				} else {
					go runB()
					synctest.Wait()
					go runA()
				}

				synctest.Wait()

				select {
				case <-doneA:
				default:
					t.Fatal("A did not complete")
				}
				select {
				case <-doneB:
				default:
					t.Fatal("B did not complete")
				}

				if state.Load()&a2Bit == 0 {
					t.Fatal("a2 did not run")
				}
				if state.Load()&b2Bit == 0 {
					t.Fatal("b2 did not run")
				}
			})
		})
	}
}

func TestRendezvousCancelWhileWaiting(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		r := rendezvous.NewRendezvous()

		var state atomic.Uint32
		exited := make(chan struct{})

		go func() {
			r.A(ctx, func() {
				state.Or(a1Bit)
			}, func() {
				state.Or(a2Bit)
			})
			close(exited)
		}()

		// Let A reach the rendezvous wait for B.
		synctest.Wait()

		cancel()
		synctest.Wait()

		select {
		case <-exited:
		default:
			t.Fatal("A did not exit after cancellation")
		}

		if state.Load()&a1Bit == 0 {
			t.Fatal("a1 did not run")
		}
		if state.Load()&a2Bit != 0 {
			t.Fatal("a2 ran despite B never arriving")
		}
	})
}
