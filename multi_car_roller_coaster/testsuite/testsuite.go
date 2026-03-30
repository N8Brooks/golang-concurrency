// Package testsuite contains reusable behavioral tests for multi-car roller coaster implementations.
package testsuite

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
)

type MultiCarRollerCoaster interface {
	Passenger(ctx context.Context, board, unboard func()) error
	Car(ctx context.Context, car int, load, run, unload func()) error
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

func startPassenger(rc MultiCarRollerCoaster, ctx context.Context) passengerRun {
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

func startCar(rc MultiCarRollerCoaster, ctx context.Context, car int) carRun {
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
		run.done <- rc.Car(ctx, car, func() {
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

func runOneRide(t *testing.T, rc MultiCarRollerCoaster, ctx context.Context, car int, carFirst bool) {
	t.Helper()

	carRun := startCar(rc, ctx, car)

	if carFirst {
		synctest.Wait()

		requireClosed(t, carRun.loaded, "car did not load")
		requireOpen(t, carRun.ran, "car ran without passengers")
		requirePending(t, carRun.done, "car returned before passengers arrived")

		p1 := startPassenger(rc, ctx)
		p2 := startPassenger(rc, ctx)
		synctest.Wait()

		requireOpen(t, p1.boarded, "passenger 1 boarded before load completed")
		requireOpen(t, p2.boarded, "passenger 2 boarded before load completed")

		close(carRun.releaseLoad)
		synctest.Wait()

		requireClosed(t, p1.boarded, "passenger 1 did not board after loading")
		requireClosed(t, p2.boarded, "passenger 2 did not board after loading")
		requireOpen(t, carRun.ran, "car ran before both passengers finished boarding")

		close(p1.releaseBoard)
		close(p2.releaseBoard)
		synctest.Wait()

		requireClosed(t, carRun.ran, "car did not start after both passengers boarded")

		close(carRun.releaseRun)
		synctest.Wait()

		requireClosed(t, carRun.unloaded, "car did not unload")
		requireOpen(t, p1.unboarded, "passenger 1 unboarded before unload completed")
		requireOpen(t, p2.unboarded, "passenger 2 unboarded before unload completed")

		close(carRun.releaseUnload)
		synctest.Wait()

		requireClosed(t, p1.unboarded, "passenger 1 did not unboard after unload")
		requireClosed(t, p2.unboarded, "passenger 2 did not unboard after unload")

		close(p1.releaseUnboard)
		close(p2.releaseUnboard)
		synctest.Wait()

		requireSuccess(t, p1.done, "passenger 1 did not complete")
		requireSuccess(t, p2.done, "passenger 2 did not complete")
		requireSuccess(t, carRun.done, "car did not complete")
		return
	}

	p1 := startPassenger(rc, ctx)
	p2 := startPassenger(rc, ctx)
	synctest.Wait()

	requireOpen(t, p1.boarded, "passenger 1 boarded before the car loaded")
	requireOpen(t, p2.boarded, "passenger 2 boarded before the car loaded")

	close(carRun.releaseLoad)
	synctest.Wait()

	requireClosed(t, p1.boarded, "passenger 1 did not board after loading")
	requireClosed(t, p2.boarded, "passenger 2 did not board after loading")
	requireOpen(t, carRun.ran, "car ran before both passengers finished boarding")

	close(p1.releaseBoard)
	close(p2.releaseBoard)
	synctest.Wait()

	requireClosed(t, carRun.ran, "car did not start after both passengers boarded")

	close(carRun.releaseRun)
	synctest.Wait()

	requireClosed(t, carRun.unloaded, "car did not unload")
	requireOpen(t, p1.unboarded, "passenger 1 unboarded before unload completed")
	requireOpen(t, p2.unboarded, "passenger 2 unboarded before unload completed")

	close(carRun.releaseUnload)
	synctest.Wait()

	requireClosed(t, p1.unboarded, "passenger 1 did not unboard after unload")
	requireClosed(t, p2.unboarded, "passenger 2 did not unboard after unload")

	close(p1.releaseUnboard)
	close(p2.releaseUnboard)
	synctest.Wait()

	requireSuccess(t, p1.done, "passenger 1 did not complete")
	requireSuccess(t, p2.done, "passenger 2 did not complete")
	requireSuccess(t, carRun.done, "car did not complete")
}

func Run(t *testing.T, newImpl func(cars, capacity int) MultiCarRollerCoaster) {
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
					runOneRide(t, newImpl(1, 2), t.Context(), 0, tc.carFirst)
				})
			})
		}
	})

	t.Run("ExclusiveBoardingAndOrderedUnloading", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(2, 2)
			ctx := t.Context()

			car0 := startCar(rc, ctx, 0)
			synctest.Wait()
			requireClosed(t, car0.loaded, "car 0 did not enter the loading area first")

			car1 := startCar(rc, ctx, 1)
			synctest.Wait()
			requireOpen(t, car1.loaded, "car 1 loaded while car 0 still had the boarding area")

			p1 := startPassenger(rc, ctx)
			p2 := startPassenger(rc, ctx)
			synctest.Wait()

			close(car0.releaseLoad)
			synctest.Wait()

			requireClosed(t, p1.boarded, "car 0 passenger 1 did not board")
			requireClosed(t, p2.boarded, "car 0 passenger 2 did not board")
			requireOpen(t, car0.ran, "car 0 ran before its passengers finished boarding")

			close(p1.releaseBoard)
			close(p2.releaseBoard)
			synctest.Wait()

			requireClosed(t, car0.ran, "car 0 did not start running")
			requireClosed(t, car1.loaded, "car 1 did not gain the loading area after car 0 boarded")

			p3 := startPassenger(rc, ctx)
			p4 := startPassenger(rc, ctx)
			synctest.Wait()

			close(car1.releaseLoad)
			synctest.Wait()

			requireClosed(t, p3.boarded, "car 1 passenger 1 did not board")
			requireClosed(t, p4.boarded, "car 1 passenger 2 did not board")

			close(p3.releaseBoard)
			close(p4.releaseBoard)
			synctest.Wait()

			requireClosed(t, car1.ran, "car 1 did not start running while car 0 was already on the track")
			close(car1.releaseRun)
			synctest.Wait()
			requireOpen(t, car1.unloaded, "car 1 unloaded before car 0")

			close(car0.releaseRun)
			synctest.Wait()

			requireClosed(t, car0.unloaded, "car 0 did not unload first")
			requireOpen(t, car1.unloaded, "car 1 unloaded before car 0 finished unloading")

			close(car0.releaseUnload)
			synctest.Wait()

			requireClosed(t, p1.unboarded, "car 0 passenger 1 did not begin unboarding")
			requireClosed(t, p2.unboarded, "car 0 passenger 2 did not begin unboarding")
			requireOpen(t, car1.unloaded, "car 1 unloaded before car 0 passengers finished unboarding")

			close(p1.releaseUnboard)
			close(p2.releaseUnboard)
			synctest.Wait()

			requireSuccess(t, car0.done, "car 0 did not complete")
			requireClosed(t, car1.unloaded, "car 1 did not unload after car 0 passengers finished unboarding")

			close(car1.releaseUnload)
			synctest.Wait()

			requireClosed(t, p3.unboarded, "car 1 passenger 1 did not begin unboarding")
			requireClosed(t, p4.unboarded, "car 1 passenger 2 did not begin unboarding")

			close(p3.releaseUnboard)
			close(p4.releaseUnboard)
			synctest.Wait()

			requireSuccess(t, p1.done, "car 0 passenger 1 did not complete")
			requireSuccess(t, p2.done, "car 0 passenger 2 did not complete")
			requireSuccess(t, p3.done, "car 1 passenger 1 did not complete")
			requireSuccess(t, p4.done, "car 1 passenger 2 did not complete")
			requireSuccess(t, car1.done, "car 1 did not complete")
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(2, 2)
			ctx := t.Context()

			for i := range 4 {
				runOneRide(t, rc, ctx, i%2, i%2 == 0)
			}
		})
	})

	t.Run("CancelWaitingPassenger", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(1, 2)
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

			replacement1 := startPassenger(rc, t.Context())
			replacement2 := startPassenger(rc, t.Context())
			car := startCar(rc, t.Context(), 0)
			synctest.Wait()

			close(car.releaseLoad)
			synctest.Wait()

			requireClosed(t, replacement1.boarded, "replacement passenger 1 did not board")
			requireClosed(t, replacement2.boarded, "replacement passenger 2 did not board")

			close(replacement1.releaseBoard)
			close(replacement2.releaseBoard)
			synctest.Wait()

			requireClosed(t, car.ran, "car did not run after cancellation freed the queue")

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

	t.Run("CancelCarWaitingForPassengers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rc := newImpl(1, 2)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			car := startCar(rc, ctx, 0)
			synctest.Wait()

			requireClosed(t, car.loaded, "car did not load")

			close(car.releaseLoad)
			synctest.Wait()

			requireOpen(t, car.ran, "car ran without enough passengers")
			requirePending(t, car.done, "car returned before cancellation")

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

			replacement1 := startPassenger(rc, t.Context())
			replacement2 := startPassenger(rc, t.Context())
			nextCar := startCar(rc, t.Context(), 0)
			synctest.Wait()

			close(nextCar.releaseLoad)
			synctest.Wait()

			requireClosed(t, replacement1.boarded, "replacement passenger 1 did not board after canceled car released the turn")
			requireClosed(t, replacement2.boarded, "replacement passenger 2 did not board after canceled car released the turn")

			close(replacement1.releaseBoard)
			close(replacement2.releaseBoard)
			synctest.Wait()
			close(nextCar.releaseRun)
			synctest.Wait()
			close(nextCar.releaseUnload)
			synctest.Wait()
			close(replacement1.releaseUnboard)
			close(replacement2.releaseUnboard)
			synctest.Wait()

			requireSuccess(t, replacement1.done, "replacement passenger 1 did not complete")
			requireSuccess(t, replacement2.done, "replacement passenger 2 did not complete")
			requireSuccess(t, nextCar.done, "replacement car did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(cars, capacity int) MultiCarRollerCoaster) {
	b.Helper()
	b.ReportAllocs()

	b.Run("Cars4Capacity8", func(b *testing.B) {
		rc := newImpl(4, 8)
		ctx, cancel := context.WithCancel(b.Context())
		defer cancel()

		carDone := make(chan error, 4)
		for i := range 4 {
			carID := i
			go func() {
				for {
					if err := rc.Car(ctx, carID, func() {}, func() {}, func() {}); err != nil {
						carDone <- err
						return
					}
				}
			}()
		}

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := rc.Passenger(ctx, func() {}, func() {}); err != nil {
					b.Fatalf("Passenger returned %v", err)
				}
			}
		})

		cancel()

		for range 4 {
			err := <-carDone
			if !errors.Is(err, context.Canceled) {
				b.Fatalf("car returned %v, want context.Canceled", err)
			}
		}
	})
}
