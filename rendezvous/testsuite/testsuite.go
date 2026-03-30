// Package testsuite contains reusable behavioral tests for rendezvous implementations.
package testsuite

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const (
	a1Bit uint32 = 1 << iota
	b1Bit
	a2Bit
	b2Bit
)

type Rendezvous interface {
	A(ctx context.Context, a1, a2 func())
	B(ctx context.Context, b1, b2 func())
}

func runOneShot(t *testing.T, r Rendezvous, ctx context.Context, aFirst bool) {
	t.Helper()

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

	if aFirst {
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
}

func Run(t *testing.T, newImpl func() Rendezvous) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name   string
			aFirst bool
		}{
			{name: "AFirst", aFirst: true},
			{name: "BFirst", aFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					r := newImpl()
					runOneShot(t, r, t.Context(), tc.aFirst)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		for _, tc := range []struct {
			name   string
			aFirst bool
		}{
			{name: "AFirst", aFirst: true},
			{name: "BFirst", aFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					r := newImpl()
					ctx := t.Context()
					for range 5 {
						runOneShot(t, r, ctx, tc.aFirst)
					}
				})
			})
		}
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		for _, tc := range []struct {
			name   string
			aFirst bool
		}{
			{name: "AFirst", aFirst: true},
			{name: "BFirst", aFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()

					r := newImpl()

					var state atomic.Uint32
					exited := make(chan struct{})

					if tc.aFirst {
						go func() {
							r.A(ctx, func() {
								state.Or(a1Bit)
							}, func() {
								state.Or(a2Bit)
							})
							close(exited)
						}()
					} else {
						go func() {
							r.B(ctx, func() {
								state.Or(b1Bit)
							}, func() {
								state.Or(b2Bit)
							})
							close(exited)
						}()
					}

					synctest.Wait()

					cancel()
					synctest.Wait()

					select {
					case <-exited:
					default:
						if tc.aFirst {
							t.Fatal("A did not exit after cancellation")
						}
						t.Fatal("B did not exit after cancellation")
					}

					if tc.aFirst {
						if state.Load()&a1Bit == 0 {
							t.Fatal("a1 did not run")
						}
						if state.Load()&a2Bit != 0 {
							t.Fatal("a2 ran despite B never arriving")
						}
					} else {
						if state.Load()&b1Bit == 0 {
							t.Fatal("b1 did not run")
						}
						if state.Load()&b2Bit != 0 {
							t.Fatal("b2 ran despite A never arriving")
						}
					}
				})
			})
		}
	})
}

func Benchmark(b *testing.B, newImpl func() Rendezvous) {
	b.Helper()
	b.ReportAllocs()

	b.Run("AFirst", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			r := newImpl()
			ready := make(chan struct{})
			done := make(chan struct{})

			go func() {
				r.A(ctx, func() {
					close(ready)
				}, func() {})
				close(done)
			}()

			<-ready
			r.B(ctx, func() {}, func() {})
			<-done
		}
	})

	b.Run("BFirst", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			r := newImpl()
			ready := make(chan struct{})
			done := make(chan struct{})

			go func() {
				r.B(ctx, func() {
					close(ready)
				}, func() {})
				close(done)
			}()

			<-ready
			r.A(ctx, func() {}, func() {})
			<-done
		}
	})
}
