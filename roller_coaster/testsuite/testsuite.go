// Package testsuite contains reusable behavioral tests for roller coaster implementations.
package testsuite

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
)

type RollerCoaster interface {
	Passenger(ctx context.Context, board, unboard func()) error
	Car(ctx context.Context, load, run, unload func()) error
}

type passengerRun struct {
	boarded        chan struct{}
	unboarded      chan struct{}
	releaseBoard   chan struct{}
	releaseUnboard chan struct{}
	done           chan error
}

type carRun struct {
	loaded        chan struct{}
	ran           chan struct{}
	unloaded      chan struct{}
	releaseLoad   chan struct{}
	releaseRun    chan struct{}
	releaseUnload chan struct{}
	done          chan error
}

func startPassenger(rc RollerCoaster, ctx context.Context) passengerRun {
	run := passengerRun{
		boarded:        make(chan struct{}),
		unboarded:      make(chan struct{}),
		releaseBoard:   make(chan struct{}),
		releaseUnboard: make(chan struct{}),
		done:           make(chan error, 1),
	}

	go func() {
		run.done <- rc.Passenger(ctx, func() {
			close(run.boarded)
			<-run.releaseBoard
		}, func() {
			close(run.unboarded)
			<-run.releaseUnboard
		})
	}()

	return run
}

func startCar(rc RollerCoaster, ctx context.Context) carRun {
	run := carRun{
		loaded:        make(chan struct{}),
		ran:           make(chan struct{}),
		unloaded:      make(chan struct{}),
		releaseLoad:   make(chan struct{}),
		releaseRun:    make(chan struct{}),
		releaseUnload: make(chan struct{}),
		done:          make(chan error, 1),
	}

	go func() {
		run.done <- rc.Car(ctx, func() {
			close(run.loaded)
			<-run.releaseLoad
		}, func() {
			close(run.ran)
			<-run.releaseRun
		}, func() {
			close(run.unloaded)
			<-run.releaseUnload
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

func countBoarded(passengers ...passengerRun) int {
	var count int
	for _, passenger := range passengers {
		select {
		case <-passenger.boarded:
			count++
		default:
		}
	}
	return count
}

func runOneRide(t *testing.T, rc RollerCoaster, ctx context.Context, carFirst bool) {
	t.Helper()

	car := startCar(rc, ctx)

	if !carFirst {
		p1 := startPassenger(rc, ctx)
		p2 := startPassenger(rc, ctx)
		synctest.Wait()

		requireOpen(t, p1.boarded, "passenger 1 boarded before the car loaded")
		requireOpen(t, p2.boarded, "passenger 2 boarded before the car loaded")
		requirePending(t, car.done, "car returned before completing the ride")

		close(car.releaseLoad)
		synctest.Wait()

		requireClosed(t, p1.boarded, "passenger 1 did not board after loading")
		requireClosed(t, p2.boarded, "passenger 2 did not board after loading")
		requireOpen(t, car.ran, "car ran before both passengers finished boarding")

		close(p1.releaseBoard)
		close(p2.releaseBoard)
		synctest.Wait()

		requireClosed(t, car.ran, "car did not start after both passengers boarded")
		requireOpen(t, p1.unboarded, "passenger 1 unboarded before the car unloaded")
		requireOpen(t, p2.unboarded, "passenger 2 unboarded before the car unloaded")

		close(car.releaseRun)
		synctest.Wait()

		requireClosed(t, car.unloaded, "car did not unload after the run")
		requireOpen(t, p1.unboarded, "passenger 1 unboarded before unload completed")
		requireOpen(t, p2.unboarded, "passenger 2 unboarded before unload completed")

		close(car.releaseUnload)
		synctest.Wait()

		requireClosed(t, p1.unboarded, "passenger 1 did not unboard after unload")
		requireClosed(t, p2.unboarded, "passenger 2 did not unboard after unload")
		requirePending(t, car.done, "car completed before both passengers finished unboarding")

		close(p1.releaseUnboard)
		close(p2.releaseUnboard)
		synctest.Wait()

		requireSuccess(t, p1.done, "passenger 1 did not complete")
		requireSuccess(t, p2.done, "passenger 2 did not complete")
		requireSuccess(t, car.done, "car did not complete")
		return
	}

	synctest.Wait()

	requireClosed(t, car.loaded, "car did not load")
	requireOpen(t, car.ran, "car ran without passengers")
	requirePending(t, car.done, "car returned before passengers arrived")

	p1 := startPassenger(rc, ctx)
	p2 := startPassenger(rc, ctx)
	synctest.Wait()

	requireOpen(t, p1.boarded, "passenger 1 boarded before load completed")
	requireOpen(t, p2.boarded, "passenger 2 boarded before load completed")

	close(car.releaseLoad)
	synctest.Wait()

	requireClosed(t, p1.boarded, "passenger 1 did not board after load completed")
	requireClosed(t, p2.boarded, "passenger 2 did not board after load completed")
	requireOpen(t, car.ran, "car ran before both passengers finished boarding")

	close(p1.releaseBoard)
	close(p2.releaseBoard)
	synctest.Wait()

	requireClosed(t, car.ran, "car did not start after both passengers boarded")
	requireOpen(t, p1.unboarded, "passenger 1 unboarded before the car unloaded")
	requireOpen(t, p2.unboarded, "passenger 2 unboarded before the car unloaded")

	close(car.releaseRun)
	synctest.Wait()

	requireClosed(t, car.unloaded, "car did not unload after the run")
	requireOpen(t, p1.unboarded, "passenger 1 unboarded before unload completed")
	requireOpen(t, p2.unboarded, "passenger 2 unboarded before unload completed")

	close(car.releaseUnload)
	synctest.Wait()

	requireClosed(t, p1.unboarded, "passenger 1 did not unboard after unload")
	requireClosed(t, p2.unboarded, "passenger 2 did not unboard after unload")
	requirePending(t, car.done, "car completed before both passengers finished unboarding")

	close(p1.releaseUnboard)
	close(p2.releaseUnboard)
	synctest.Wait()

	requireSuccess(t, p1.done, "passenger 1 did not complete")
	requireSuccess(t, p2.done, "passenger 2 did not complete")
	requireSuccess(t, car.done, "car did not complete")
}

func Run(t *testing.T, newImpl func(capacity int) RollerCoaster) {
	t.Helper()

	t.Run("OneRide", func(t *testing.T) {
		for _, tc := range []struct {
			name     string
			carFirst bool
		}{
			{name: "PassengersFirst", carFirst: false},
			{name: "CarFirst", carFirst: true},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneRide(t, newImpl(2), t.Context(), tc.carFirst)
				})
			})
		}
	})

	t.Run("BoardsExactlyCapacity", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(2)
			ctx := t.Context()

			p1 := startPassenger(rc, ctx)
			p2 := startPassenger(rc, ctx)
			extraCtx, cancelExtra := context.WithCancel(ctx)
			defer cancelExtra()
			p3 := startPassenger(rc, extraCtx)
			car := startCar(rc, ctx)

			synctest.Wait()

			requireOpen(t, p1.boarded, "passenger 1 boarded before the car finished loading")
			requireOpen(t, p2.boarded, "passenger 2 boarded before the car finished loading")
			requireOpen(t, p3.boarded, "extra passenger boarded before the car finished loading")

			close(car.releaseLoad)
			synctest.Wait()

			if got := countBoarded(p1, p2, p3); got != 2 {
				t.Fatalf("got %d boarded passengers for a capacity-2 car, want 2", got)
			}
			requireOpen(t, p3.unboarded, "extra passenger unboarded without boarding")
			requireOpen(t, car.ran, "car ran before both selected passengers finished boarding")

			switch {
			case isClosed(p1.boarded) && isClosed(p2.boarded):
				close(p1.releaseBoard)
				close(p2.releaseBoard)
			case isClosed(p1.boarded) && isClosed(p3.boarded):
				close(p1.releaseBoard)
				close(p3.releaseBoard)
			case isClosed(p2.boarded) && isClosed(p3.boarded):
				close(p2.releaseBoard)
				close(p3.releaseBoard)
			default:
				t.Fatal("no passengers boarded")
			}

			synctest.Wait()

			requireClosed(t, car.ran, "car did not start after its selected passengers boarded")

			close(car.releaseRun)
			synctest.Wait()

			requireClosed(t, car.unloaded, "car did not unload")

			for _, passenger := range []passengerRun{p1, p2, p3} {
				if isClosed(passenger.boarded) {
					requireOpen(t, passenger.unboarded, "boarded passenger unboarded before unload completed")
				}
			}

			close(car.releaseUnload)
			synctest.Wait()

			for _, passenger := range []passengerRun{p1, p2, p3} {
				if isClosed(passenger.unboarded) {
					close(passenger.releaseUnboard)
				} else if isClosed(passenger.boarded) {
					t.Fatal("boarded passenger did not unboard after unload completed")
				}
			}

			synctest.Wait()

			for _, passenger := range []passengerRun{p1, p2, p3} {
				if isClosed(passenger.unboarded) {
					requireSuccess(t, passenger.done, "selected passenger did not complete")
				}
			}
			requireSuccess(t, car.done, "car did not complete")

			cancelExtra()
			synctest.Wait()

			if !isClosed(p3.boarded) {
				select {
				case err := <-p3.done:
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("extra passenger returned %v, want context.Canceled", err)
					}
				default:
					t.Fatal("extra passenger did not exit after cancellation")
				}
			}
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(2)
			ctx := t.Context()

			for i := range 4 {
				runOneRide(t, rc, ctx, i%2 == 0)
			}
		})
	})

	t.Run("CancelWaitingPassenger", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(2)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			passenger := startPassenger(rc, ctx)
			synctest.Wait()

			requireOpen(t, passenger.boarded, "passenger boarded without a car")
			requirePending(t, passenger.done, "passenger returned before cancellation")

			cancel()
			synctest.Wait()

			select {
			case err := <-passenger.done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("passenger returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("passenger did not exit after cancellation")
			}

			requireOpen(t, passenger.boarded, "passenger boarded despite being canceled while waiting")
			requireOpen(t, passenger.unboarded, "passenger unboarded despite never boarding")

			replacement1 := startPassenger(rc, t.Context())
			replacement2 := startPassenger(rc, t.Context())
			car := startCar(rc, t.Context())
			synctest.Wait()

			close(car.releaseLoad)
			synctest.Wait()

			requireClosed(t, replacement1.boarded, "replacement passenger 1 did not board")
			requireClosed(t, replacement2.boarded, "replacement passenger 2 did not board")

			close(replacement1.releaseBoard)
			close(replacement2.releaseBoard)
			synctest.Wait()

			requireClosed(t, car.ran, "car did not run with replacement passengers")

			close(car.releaseRun)
			synctest.Wait()
			close(car.releaseUnload)
			synctest.Wait()

			close(replacement1.releaseUnboard)
			close(replacement2.releaseUnboard)
			synctest.Wait()

			requireSuccess(t, replacement1.done, "replacement passenger 1 did not complete")
			requireSuccess(t, replacement2.done, "replacement passenger 2 did not complete")
			requireSuccess(t, car.done, "car did not complete after replacement ride")
		})
	})

	t.Run("CancelWaitingCar", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(2)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			car := startCar(rc, ctx)
			synctest.Wait()

			requireClosed(t, car.loaded, "car did not load")
			requireOpen(t, car.ran, "car ran without enough passengers")
			requirePending(t, car.done, "car returned before cancellation")

			close(car.releaseLoad)
			synctest.Wait()

			requireOpen(t, car.ran, "car ran without enough passengers")
			requirePending(t, car.done, "car returned before cancellation while waiting for passengers")

			cancel()
			synctest.Wait()

			select {
			case err := <-car.done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("car returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("car did not exit after cancellation")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) RollerCoaster) {
	b.Helper()
	b.ReportAllocs()

	b.Run("Capacity8", func(b *testing.B) {
		rc := newImpl(8)
		ctx, cancel := context.WithCancel(b.Context())
		defer cancel()

		carDone := make(chan error, 1)
		go func() {
			for {
				if err := rc.Car(ctx, func() {}, func() {}, func() {}); err != nil {
					carDone <- err
					return
				}
			}
		}()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := rc.Passenger(ctx, func() {}, func() {}); err != nil {
					b.Fatalf("Passenger returned %v", err)
				}
			}
		})

		cancel()

		err := <-carDone
		if !errors.Is(err, context.Canceled) {
			b.Fatalf("car returned %v, want context.Canceled", err)
		}
	})
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
