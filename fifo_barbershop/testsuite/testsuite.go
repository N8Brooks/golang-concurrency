// Package testsuite contains reusable behavioral tests for FIFO barbershop implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
)

type FIFOBarbershop interface {
	Customer(ctx context.Context, getHairCut, balk func()) error
	Barber(ctx context.Context, cutHair func()) error
}

type customerRun struct {
	haircutStarted chan struct{}
	balked         chan struct{}
	release        chan struct{}
	done           chan error
}

type barberRun struct {
	started chan struct{}
	release chan struct{}
	done    chan error
}

func startCustomer(shop FIFOBarbershop, ctx context.Context) customerRun {
	run := customerRun{
		haircutStarted: make(chan struct{}),
		balked:         make(chan struct{}),
		release:        make(chan struct{}),
		done:           make(chan error, 1),
	}

	go func() {
		run.done <- shop.Customer(ctx, func() {
			close(run.haircutStarted)
			<-run.release
		}, func() {
			close(run.balked)
		})
	}()

	return run
}

func startBarber(shop FIFOBarbershop, ctx context.Context) barberRun {
	run := barberRun{
		started: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		run.done <- shop.Barber(ctx, func() {
			close(run.started)
			<-run.release
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

func countStarted(customers ...customerRun) (count int, index int) {
	index = -1
	for i, customer := range customers {
		select {
		case <-customer.haircutStarted:
			count++
			index = i
		default:
		}
	}
	return count, index
}

func runOneShot(t *testing.T, shop FIFOBarbershop, ctx context.Context, barberFirst bool) {
	t.Helper()

	var barber barberRun
	var customer customerRun

	if barberFirst {
		barber = startBarber(shop, ctx)
		synctest.Wait()

		requireOpen(t, barber.started, "barber started cutting hair without a customer")
		requirePending(t, barber.done, "barber returned before a customer arrived")

		customer = startCustomer(shop, ctx)
	} else {
		customer = startCustomer(shop, ctx)
		synctest.Wait()

		requireOpen(t, customer.haircutStarted, "customer started haircut without a barber")
		requireOpen(t, customer.balked, "customer balked despite available capacity")
		requirePending(t, customer.done, "customer returned before being matched")

		barber = startBarber(shop, ctx)
	}

	synctest.Wait()

	requireClosed(t, barber.started, "barber did not start cutting hair")
	requireClosed(t, customer.haircutStarted, "customer did not start haircut")
	requireOpen(t, customer.balked, "matched customer balked")

	close(customer.release)
	close(barber.release)
	synctest.Wait()

	requireSuccess(t, customer.done, "customer did not complete")
	requireSuccess(t, barber.done, "barber did not complete")
}

func Run(t *testing.T, newImpl func(capacity int) FIFOBarbershop) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name        string
			barberFirst bool
		}{
			{name: "BarberFirst", barberFirst: true},
			{name: "CustomerFirst", barberFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneShot(t, newImpl(2), t.Context(), tc.barberFirst)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		for _, tc := range []struct {
			name        string
			barberFirst bool
		}{
			{name: "BarberFirst", barberFirst: true},
			{name: "CustomerFirst", barberFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					shop := newImpl(3)
					ctx := t.Context()
					for range 5 {
						runOneShot(t, shop, ctx, tc.barberFirst)
					}
				})
			})
		}
	})

	t.Run("FullShopBalks", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			shop := newImpl(2)
			ctx := t.Context()

			first := startCustomer(shop, ctx)
			second := startCustomer(shop, ctx)
			synctest.Wait()

			for i, customer := range []customerRun{first, second} {
				requireOpen(t, customer.haircutStarted, fmt.Sprintf("customer %d started haircut without a barber", i))
				requireOpen(t, customer.balked, fmt.Sprintf("customer %d balked despite available capacity", i))
				requirePending(t, customer.done, fmt.Sprintf("customer %d returned before being served", i))
			}

			third := startCustomer(shop, ctx)
			synctest.Wait()

			requireClosed(t, third.balked, "customer did not balk when the shop was full")
			requireOpen(t, third.haircutStarted, "balking customer started haircut")
			requireSuccess(t, third.done, "balking customer did not return")

			barber1 := startBarber(shop, ctx)
			synctest.Wait()

			count, servedIndex := countStarted(first, second)
			if count != 1 {
				t.Fatalf("got %d customers in the barber chair with one barber, want 1", count)
			}

			served := []customerRun{first, second}[servedIndex]
			waiting := []customerRun{first, second}[1-servedIndex]

			requireClosed(t, barber1.started, "barber did not start first haircut")
			requireOpen(t, waiting.haircutStarted, "second waiting customer started too early")

			close(served.release)
			close(barber1.release)
			synctest.Wait()

			requireSuccess(t, served.done, "first served customer did not complete")
			requireSuccess(t, barber1.done, "barber did not finish first haircut")

			barber2 := startBarber(shop, ctx)
			synctest.Wait()

			requireClosed(t, barber2.started, "barber did not start second haircut")
			requireClosed(t, waiting.haircutStarted, "waiting customer did not get the second haircut")

			close(waiting.release)
			close(barber2.release)
			synctest.Wait()

			requireSuccess(t, waiting.done, "second served customer did not complete")
			requireSuccess(t, barber2.done, "barber did not finish second haircut")
		})
	})

	t.Run("OneCustomerPerHaircut", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			shop := newImpl(2)
			ctx := t.Context()

			first := startCustomer(shop, ctx)
			second := startCustomer(shop, ctx)
			synctest.Wait()

			barberDone := make(chan error, 2)
			firstCutStarted := make(chan struct{})
			firstCutRelease := make(chan struct{})
			secondCutStarted := make(chan struct{})
			secondCutRelease := make(chan struct{})

			go func() {
				barberDone <- shop.Barber(ctx, func() {
					close(firstCutStarted)
					<-firstCutRelease
				})
				barberDone <- shop.Barber(ctx, func() {
					close(secondCutStarted)
					<-secondCutRelease
				})
			}()

			synctest.Wait()

			requireClosed(t, firstCutStarted, "barber did not start the first haircut")

			count, servedIndex := countStarted(first, second)
			if count != 1 {
				t.Fatalf("got %d customers in the barber chair during one haircut, want 1", count)
			}

			served := []customerRun{first, second}[servedIndex]
			waiting := []customerRun{first, second}[1-servedIndex]

			requireOpen(t, waiting.haircutStarted, "second customer started haircut before the first haircut completed")
			requireOpen(t, secondCutStarted, "barber started a second haircut before the first finished")

			close(served.release)
			close(firstCutRelease)
			synctest.Wait()

			requireClosed(t, secondCutStarted, "barber did not start the second haircut")
			requireClosed(t, waiting.haircutStarted, "waiting customer did not start the second haircut")

			close(waiting.release)
			close(secondCutRelease)
			synctest.Wait()

			requireSuccess(t, first.done, "first customer did not complete")
			requireSuccess(t, second.done, "second customer did not complete")
			requireSuccess(t, barberDone, "barber did not complete first haircut")
			requireSuccess(t, barberDone, "barber did not complete second haircut")
		})
	})

	t.Run("ServesCustomersFIFO", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			shop := newImpl(3)
			ctx := t.Context()

			first := startCustomer(shop, ctx)
			synctest.Wait()
			requireOpen(t, first.haircutStarted, "first customer started haircut without a barber")
			requireOpen(t, first.balked, "first customer balked despite available capacity")

			second := startCustomer(shop, ctx)
			synctest.Wait()
			requireOpen(t, second.haircutStarted, "second customer started haircut without a barber")
			requireOpen(t, second.balked, "second customer balked despite available capacity")

			third := startCustomer(shop, ctx)
			synctest.Wait()
			requireOpen(t, third.haircutStarted, "third customer started haircut without a barber")
			requireOpen(t, third.balked, "third customer balked despite available capacity")

			firstCutStarted := make(chan struct{})
			firstCutRelease := make(chan struct{})
			secondCutStarted := make(chan struct{})
			secondCutRelease := make(chan struct{})
			thirdCutStarted := make(chan struct{})
			thirdCutRelease := make(chan struct{})
			barberDone := make(chan error, 3)

			go func() {
				barberDone <- shop.Barber(ctx, func() {
					close(firstCutStarted)
					<-firstCutRelease
				})
				barberDone <- shop.Barber(ctx, func() {
					close(secondCutStarted)
					<-secondCutRelease
				})
				barberDone <- shop.Barber(ctx, func() {
					close(thirdCutStarted)
					<-thirdCutRelease
				})
			}()

			synctest.Wait()

			requireClosed(t, firstCutStarted, "barber did not start the first haircut")
			requireClosed(t, first.haircutStarted, "first customer was not served first")
			requireOpen(t, second.haircutStarted, "second customer was served before the first haircut finished")
			requireOpen(t, third.haircutStarted, "third customer was served before the first haircut finished")
			requireOpen(t, secondCutStarted, "barber started the second haircut before the first finished")
			requireOpen(t, thirdCutStarted, "barber started the third haircut before the first finished")

			close(first.release)
			close(firstCutRelease)
			synctest.Wait()

			requireClosed(t, secondCutStarted, "barber did not start the second haircut")
			requireClosed(t, second.haircutStarted, "second customer was not served second")
			requireOpen(t, third.haircutStarted, "third customer was served before the second haircut finished")
			requireOpen(t, thirdCutStarted, "barber started the third haircut before the second finished")

			close(second.release)
			close(secondCutRelease)
			synctest.Wait()

			requireClosed(t, thirdCutStarted, "barber did not start the third haircut")
			requireClosed(t, third.haircutStarted, "third customer was not served third")

			close(third.release)
			close(thirdCutRelease)
			synctest.Wait()

			requireSuccess(t, first.done, "first customer did not complete")
			requireSuccess(t, second.done, "second customer did not complete")
			requireSuccess(t, third.done, "third customer did not complete")
			requireSuccess(t, barberDone, "barber did not complete first haircut")
			requireSuccess(t, barberDone, "barber did not complete second haircut")
			requireSuccess(t, barberDone, "barber did not complete third haircut")
		})
	})

	t.Run("CancelWaitingCustomer", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			shop := newImpl(1)

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			customer := startCustomer(shop, ctx)
			synctest.Wait()

			requireOpen(t, customer.haircutStarted, "customer started haircut without a barber")
			requireOpen(t, customer.balked, "customer balked despite entering the shop")
			requirePending(t, customer.done, "customer returned before cancellation")

			cancel()
			synctest.Wait()

			select {
			case err := <-customer.done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("customer returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("customer did not exit after cancellation")
			}

			requireOpen(t, customer.haircutStarted, "canceled customer still got a haircut")
			requireOpen(t, customer.balked, "canceled customer balked instead of canceling")

			replacement := startCustomer(shop, t.Context())
			barber := startBarber(shop, t.Context())
			synctest.Wait()

			requireClosed(t, replacement.haircutStarted, "replacement customer did not get served after cancellation freed a spot")
			requireClosed(t, barber.started, "barber did not serve replacement customer")

			close(replacement.release)
			close(barber.release)
			synctest.Wait()

			requireSuccess(t, replacement.done, "replacement customer did not complete")
			requireSuccess(t, barber.done, "barber did not complete replacement haircut")
		})
	})

	t.Run("CancelSleepingBarber", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			shop := newImpl(2)

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			barber := startBarber(shop, ctx)
			synctest.Wait()

			requireOpen(t, barber.started, "barber started cutting hair without a customer")
			requirePending(t, barber.done, "barber returned before cancellation")

			cancel()
			synctest.Wait()

			select {
			case err := <-barber.done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("barber returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("barber did not exit after cancellation")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) FIFOBarbershop) {
	b.Helper()
	b.ReportAllocs()

	b.Run("Capacity16", func(b *testing.B) {
		shop := newImpl(16)
		ctx, cancel := context.WithCancel(b.Context())
		defer cancel()

		barberDone := make(chan error, 1)
		go func() {
			for {
				if err := shop.Barber(ctx, func() {}); err != nil {
					barberDone <- err
					return
				}
			}
		}()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := shop.Customer(ctx, func() {}, func() {}); err != nil {
					b.Fatalf("Customer returned %v", err)
				}
			}
		})

		cancel()

		err := <-barberDone
		if !errors.Is(err, context.Canceled) {
			b.Fatalf("barber returned %v, want context.Canceled", err)
		}
	})
}
