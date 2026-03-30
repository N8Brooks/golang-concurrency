// Package testsuite contains reusable behavioral tests for baboon crossing implementations.
package testsuite

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type Crossing interface {
	Left(ctx context.Context, arrive, cross, exit func()) error
	Right(ctx context.Context, arrive, cross, exit func()) error
}

func requireClosed(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	default:
		t.Fatal(msg)
	}
}

func runOneShot(t *testing.T, c Crossing, ctx context.Context, leftFirst bool) {
	t.Helper()

	var leftActive atomic.Int32
	var rightActive atomic.Int32

	leftRelease := make(chan struct{})
	rightRelease := make(chan struct{})
	var leftReleaseOnce sync.Once
	var rightReleaseOnce sync.Once
	defer leftReleaseOnce.Do(func() { close(leftRelease) })
	defer rightReleaseOnce.Do(func() { close(rightRelease) })

	leftDone := make(chan error, 1)
	rightDone := make(chan error, 1)

	runLeft := func() {
		leftDone <- c.Left(ctx, func() {}, func() {
			if rightActive.Load() != 0 {
				t.Error("left baboon crossed concurrently with a right baboon")
			}
			leftActive.Add(1)
			<-leftRelease
			leftActive.Add(-1)
		}, func() {})
	}
	runRight := func() {
		rightDone <- c.Right(ctx, func() {}, func() {
			if leftActive.Load() != 0 {
				t.Error("right baboon crossed concurrently with a left baboon")
			}
			rightActive.Add(1)
			<-rightRelease
			rightActive.Add(-1)
		}, func() {})
	}

	if leftFirst {
		go runLeft()
		synctest.Wait()
		go runRight()
		synctest.Wait()

		if leftActive.Load() != 1 {
			t.Fatalf("got %d active left baboons, want 1", leftActive.Load())
		}
		if rightActive.Load() != 0 {
			t.Fatalf("got %d active right baboons while left was on the rope, want 0", rightActive.Load())
		}

		leftReleaseOnce.Do(func() { close(leftRelease) })
		synctest.Wait()
		rightReleaseOnce.Do(func() { close(rightRelease) })
	} else {
		go runRight()
		synctest.Wait()
		go runLeft()
		synctest.Wait()

		if rightActive.Load() != 1 {
			t.Fatalf("got %d active right baboons, want 1", rightActive.Load())
		}
		if leftActive.Load() != 0 {
			t.Fatalf("got %d active left baboons while right was on the rope, want 0", leftActive.Load())
		}

		rightReleaseOnce.Do(func() { close(rightRelease) })
		synctest.Wait()
		leftReleaseOnce.Do(func() { close(leftRelease) })
	}

	synctest.Wait()

	select {
	case err := <-leftDone:
		if err != nil {
			t.Fatalf("left returned %v", err)
		}
	default:
		t.Fatal("left did not complete")
	}

	select {
	case err := <-rightDone:
		if err != nil {
			t.Fatalf("right returned %v", err)
		}
	default:
		t.Fatal("right did not complete")
	}
}

func Run(t *testing.T, newImpl func() Crossing) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name      string
			leftFirst bool
		}{
			{name: "LeftFirst", leftFirst: true},
			{name: "RightFirst", leftFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneShot(t, newImpl(), t.Context(), tc.leftFirst)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := newImpl()
			ctx := t.Context()
			for i := range 4 {
				runOneShot(t, c, ctx, i%2 == 0)
			}
		})
	})

	t.Run("MaxFiveSameDirection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			c := newImpl()

			var active atomic.Int32
			var maxActive atomic.Int32
			release := make(chan struct{})
			var releaseOnce sync.Once
			defer releaseOnce.Do(func() { close(release) })

			done := make([]chan error, 6)
			for i := range 6 {
				done[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- c.Left(ctx, func() {}, func() {
						current := active.Add(1)
						for {
							max := maxActive.Load()
							if current <= max || maxActive.CompareAndSwap(max, current) {
								break
							}
						}
						<-release
						active.Add(-1)
					}, func() {})
				}(done[i])
			}

			synctest.Wait()

			if active.Load() != 5 {
				t.Fatalf("got %d active baboons, want 5", active.Load())
			}
			if maxActive.Load() != 5 {
				t.Fatalf("got max %d active baboons, want 5", maxActive.Load())
			}

			releaseOnce.Do(func() { close(release) })
			synctest.Wait()

			for i, ch := range done {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("baboon %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("baboon %d did not complete", i+1)
				}
			}
		})
	})

	t.Run("NoOppositeDirectionConcurrency", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			c := newImpl()

			var leftActive atomic.Int32
			var rightActive atomic.Int32

			leftRelease := make(chan struct{})
			var leftReleaseOnce sync.Once
			defer leftReleaseOnce.Do(func() { close(leftRelease) })

			leftDone := make([]chan error, 3)
			for i := range 3 {
				leftDone[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- c.Left(ctx, func() {}, func() {
						if rightActive.Load() != 0 {
							t.Error("left baboon crossed concurrently with right baboons")
						}
						leftActive.Add(1)
						<-leftRelease
						leftActive.Add(-1)
					}, func() {})
				}(leftDone[i])
			}

			synctest.Wait()

			rightDone := make(chan error, 1)
			go func() {
				rightDone <- c.Right(ctx, func() {}, func() {
					if leftActive.Load() != 0 {
						t.Error("right baboon crossed concurrently with left baboons")
					}
					rightActive.Add(1)
					rightActive.Add(-1)
				}, func() {})
			}()

			synctest.Wait()

			if rightActive.Load() != 0 {
				t.Fatalf("got %d active right baboons while left baboons were on the rope, want 0", rightActive.Load())
			}

			leftReleaseOnce.Do(func() { close(leftRelease) })
			synctest.Wait()

			select {
			case err := <-rightDone:
				if err != nil {
					t.Fatalf("right returned %v", err)
				}
			default:
				t.Fatal("right did not complete after left baboons exited")
			}

			for i, ch := range leftDone {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("left baboon %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("left baboon %d did not complete", i+1)
				}
			}
		})
	})

	t.Run("NoStarvationOnDirectionSwitch", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			c := newImpl()

			leftRelease := make(chan struct{})
			var leftReleaseOnce sync.Once
			defer leftReleaseOnce.Do(func() { close(leftRelease) })

			var initialLeftActive atomic.Int32
			initialLeftDone := make([]chan error, 5)
			for i := range 5 {
				initialLeftDone[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- c.Left(ctx, func() {}, func() {
						initialLeftActive.Add(1)
						<-leftRelease
						initialLeftActive.Add(-1)
					}, func() {})
				}(initialLeftDone[i])
			}

			synctest.Wait()

			if initialLeftActive.Load() != 5 {
				t.Fatalf("got %d initial left baboons on the rope, want 5", initialLeftActive.Load())
			}

			var rightActive atomic.Int32
			rightDone := make(chan error, 1)
			go func() {
				rightDone <- c.Right(ctx, func() {}, func() {
					rightActive.Add(1)
					rightActive.Add(-1)
				}, func() {})
			}()

			synctest.Wait()

			lateLeftDone := make([]chan error, 2)
			var lateLeftActive atomic.Int32
			for i := range 2 {
				lateLeftDone[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- c.Left(ctx, func() {}, func() {
						lateLeftActive.Add(1)
						lateLeftActive.Add(-1)
					}, func() {})
				}(lateLeftDone[i])
			}

			synctest.Wait()

			if rightActive.Load() != 0 {
				t.Fatalf("got %d active right baboons before the initial left batch exited, want 0", rightActive.Load())
			}
			if lateLeftActive.Load() != 0 {
				t.Fatalf("got %d late left baboons on the rope before direction switch, want 0", lateLeftActive.Load())
			}

			leftReleaseOnce.Do(func() { close(leftRelease) })
			synctest.Wait()

			select {
			case err := <-rightDone:
				if err != nil {
					t.Fatalf("right returned %v", err)
				}
			default:
				t.Fatal("right did not complete after the left batch exited")
			}

			if lateLeftActive.Load() != 0 {
				t.Fatal("late left baboons overtook the waiting right baboon")
			}

			for i, ch := range initialLeftDone {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("initial left baboon %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("initial left baboon %d did not complete", i+1)
				}
			}

			for i, ch := range lateLeftDone {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("late left baboon %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("late left baboon %d did not complete", i+1)
				}
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			c := newImpl()

			leftRelease := make(chan struct{})
			var leftReleaseOnce sync.Once
			defer leftReleaseOnce.Do(func() { close(leftRelease) })

			holderDone := make([]chan error, 5)
			for i := range 5 {
				holderDone[i] = make(chan error, 1)
				go func(ch chan error) {
					ch <- c.Left(ctx, func() {}, func() {
						<-leftRelease
					}, func() {})
				}(holderDone[i])
			}

			synctest.Wait()

			waitCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			var rightArrived atomic.Bool
			var rightExited atomic.Bool
			rightDone := make(chan error, 1)
			go func() {
				rightDone <- c.Right(waitCtx, func() {
					rightArrived.Store(true)
				}, func() {
					t.Error("canceled right baboon began crossing")
				}, func() {
					rightExited.Store(true)
				})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			select {
			case err := <-rightDone:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("right returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("right did not exit after cancellation")
			}

			if !rightArrived.Load() {
				t.Fatal("right arrival callback did not run")
			}
			if rightExited.Load() {
				t.Fatal("right exit callback ran despite cancellation")
			}

			leftReleaseOnce.Do(func() { close(leftRelease) })
			synctest.Wait()

			for i, ch := range holderDone {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("left baboon %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("left baboon %d did not complete", i+1)
				}
			}

			followRightDone := make(chan error, 1)
			go func() {
				followRightDone <- c.Right(ctx, func() {}, func() {}, func() {})
			}()

			synctest.Wait()

			select {
			case err := <-followRightDone:
				if err != nil {
					t.Fatalf("follow-up right returned %v", err)
				}
			default:
				t.Fatal("follow-up right did not complete after canceled waiter was removed")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() Crossing) {
	b.Helper()
	b.ReportAllocs()

	b.Run("SameDirection", func(b *testing.B) {
		ctx := b.Context()
		c := newImpl()
		for b.Loop() {
			done := make(chan error, 1)
			go func() {
				done <- c.Left(ctx, func() {}, func() {}, func() {})
			}()
			if err := <-done; err != nil {
				b.Fatalf("left returned %v", err)
			}
		}
	})

	b.Run("DirectionSwitch", func(b *testing.B) {
		ctx := b.Context()
		c := newImpl()
		for b.Loop() {
			leftDone := make(chan error, 1)
			rightDone := make(chan error, 1)
			go func() {
				leftDone <- c.Left(ctx, func() {}, func() {}, func() {})
			}()
			go func() {
				rightDone <- c.Right(ctx, func() {}, func() {}, func() {})
			}()
			if err := <-leftDone; err != nil {
				b.Fatalf("left returned %v", err)
			}
			if err := <-rightDone; err != nil {
				b.Fatalf("right returned %v", err)
			}
		}
	})
}
