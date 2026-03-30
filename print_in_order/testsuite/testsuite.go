// Package testsuite contains reusable behavioral tests for print-in-order implementations.
package testsuite

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const (
	firstBit uint32 = 1 << iota
	secondBit
	thirdBit
)

type PrintInOrder interface {
	First(ctx context.Context, first func()) error
	Second(ctx context.Context, second func()) error
	Third(ctx context.Context, third func()) error
}

func runOneShot(t *testing.T, p PrintInOrder, ctx context.Context, order [3]int) {
	t.Helper()

	var state atomic.Uint32
	firstDone := make(chan error, 1)
	secondDone := make(chan error, 1)
	thirdDone := make(chan error, 1)

	first := func() {
		if state.Or(firstBit)&(secondBit|thirdBit) != 0 {
			t.Error("first ran after a later callback")
		}
	}
	second := func() {
		current := state.Or(secondBit)
		if current&firstBit == 0 {
			t.Error("second ran before first")
		}
		if current&thirdBit != 0 {
			t.Error("second ran after third")
		}
	}
	third := func() {
		if state.Or(thirdBit)&(firstBit|secondBit) != (firstBit | secondBit) {
			t.Error("third ran before first and second completed")
		}
	}

	start := func(step int) {
		switch step {
		case 1:
			go func() {
				firstDone <- p.First(ctx, first)
			}()
		case 2:
			go func() {
				secondDone <- p.Second(ctx, second)
			}()
		case 3:
			go func() {
				thirdDone <- p.Third(ctx, third)
			}()
		default:
			t.Fatalf("invalid step %d", step)
		}
	}

	for _, step := range order {
		start(step)
		synctest.Wait()
	}

	synctest.Wait()

	for name, done := range map[string]chan error{
		"First":  firstDone,
		"Second": secondDone,
		"Third":  thirdDone,
	} {
		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("%s returned %v", name, err)
			}
		default:
			t.Fatalf("%s did not complete", name)
		}
	}

	if state.Load() != (firstBit | secondBit | thirdBit) {
		t.Fatalf("callbacks produced state %b, want all three callbacks to run", state.Load())
	}
}

func Run(t *testing.T, newImpl func() PrintInOrder) {
	t.Helper()

	orders := [][3]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}

	t.Run("OneShot", func(t *testing.T) {
		for _, order := range orders {
			order := order
			t.Run(orderName(order), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneShot(t, newImpl(), t.Context(), order)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			p := newImpl()
			ctx := t.Context()
			for range 5 {
				for _, order := range orders {
					runOneShot(t, p, ctx, order)
				}
			}
		})
	})

	t.Run("CancelSecondWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			p := newImpl()
			var ran atomic.Bool
			done := make(chan error, 1)

			go func() {
				done <- p.Second(ctx, func() {
					ran.Store(true)
				})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("Second returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("Second did not exit after cancellation")
			}

			if ran.Load() {
				t.Fatal("second callback ran despite cancellation while waiting")
			}
		})
	})

	t.Run("CancelThirdWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			p := newImpl()
			firstDone := make(chan error, 1)
			thirdDone := make(chan error, 1)
			var firstRan atomic.Bool
			var thirdRan atomic.Bool

			go func() {
				firstDone <- p.First(t.Context(), func() {
					firstRan.Store(true)
				})
			}()
			synctest.Wait()

			go func() {
				thirdDone <- p.Third(ctx, func() {
					thirdRan.Store(true)
				})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			select {
			case err := <-firstDone:
				if err != nil {
					t.Fatalf("First returned %v", err)
				}
			default:
				t.Fatal("First did not complete")
			}

			select {
			case err := <-thirdDone:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("Third returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("Third did not exit after cancellation")
			}

			if !firstRan.Load() {
				t.Fatal("first callback did not run")
			}
			if thirdRan.Load() {
				t.Fatal("third callback ran despite cancellation while waiting")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() PrintInOrder) {
	b.Helper()
	b.ReportAllocs()

	b.Run("Ordered", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			p := newImpl()
			if err := p.First(ctx, func() {}); err != nil {
				b.Fatalf("First returned %v", err)
			}
			if err := p.Second(ctx, func() {}); err != nil {
				b.Fatalf("Second returned %v", err)
			}
			if err := p.Third(ctx, func() {}); err != nil {
				b.Fatalf("Third returned %v", err)
			}
		}
	})

	b.Run("Reverse", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			p := newImpl()
			secondDone := make(chan error, 1)
			thirdDone := make(chan error, 1)

			go func() {
				secondDone <- p.Second(ctx, func() {})
			}()
			go func() {
				thirdDone <- p.Third(ctx, func() {})
			}()

			if err := p.First(ctx, func() {}); err != nil {
				b.Fatalf("First returned %v", err)
			}

			select {
			case err := <-secondDone:
				if err != nil {
					b.Fatalf("Second returned %v", err)
				}
			case <-ctx.Done():
				b.Fatal(ctx.Err())
			}

			select {
			case err := <-thirdDone:
				if err != nil {
					b.Fatalf("Third returned %v", err)
				}
			case <-ctx.Done():
				b.Fatal(ctx.Err())
			}
		}
	})
}

func orderName(order [3]int) string {
	name := [3]byte{'0' + byte(order[0]), '0' + byte(order[1]), '0' + byte(order[2])}
	return string(name[:])
}
