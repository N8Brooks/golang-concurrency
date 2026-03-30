// Package testsuite contains reusable behavioral tests for Hilzer's barbershop implementations.
package testsuite

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type Barbershop interface {
	Customer(ctx context.Context, enterShop, sitOnSofa, getHairCut, pay, exitShop, balk func()) error
	Barber(ctx context.Context, cutHair, acceptPayment func()) error
}

func requireClosed(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	default:
		t.Fatal(msg)
	}
}

func startBarber(ctx context.Context, shop Barbershop, cutHair, acceptPayment func()) chan struct{} {
	done := make(chan struct{})
	go func() {
		_ = shop.Barber(ctx, cutHair, acceptPayment)
		close(done)
	}()
	return done
}

func startCustomer(ctx context.Context, shop Barbershop, enterShop, sitOnSofa, getHairCut, pay, exitShop, balk func()) chan error {
	done := make(chan error, 1)
	go func() {
		done <- shop.Customer(ctx, enterShop, sitOnSofa, getHairCut, pay, exitShop, balk)
	}()
	return done
}

func Run(t *testing.T, newImpl func() Barbershop) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			shop := newImpl()

			var state atomic.Uint32
			const (
				enteredBit uint32 = 1 << iota
				sofaBit
				customerCutBit
				barberCutBit
				payBit
				acceptBit
				exitBit
			)

			releaseCut := make(chan struct{})
			var releaseCutOnce sync.Once
			defer releaseCutOnce.Do(func() { close(releaseCut) })
			barberDone := startBarber(ctx, shop, func() {
				if state.Or(barberCutBit)&sofaBit == 0 {
					t.Error("barber started cutting before the customer sat on the sofa")
				}
				<-releaseCut
			}, func() {
				if state.Or(acceptBit)&payBit == 0 {
					t.Error("barber accepted payment before the customer paid")
				}
			})

			customerDone := startCustomer(ctx, shop, func() {
				state.Or(enteredBit)
			}, func() {
				if state.Or(sofaBit)&enteredBit == 0 {
					t.Error("customer sat on the sofa before entering the shop")
				}
			}, func() {
				if state.Or(customerCutBit)&sofaBit == 0 {
					t.Error("customer started haircut before sitting on the sofa")
				}
			}, func() {
				if state.Load()&barberCutBit == 0 {
					t.Error("customer paid before the barber finished cutting hair")
				}
				state.Or(payBit)
			}, func() {
				if state.Or(exitBit)&acceptBit == 0 {
					t.Error("customer exited before payment was accepted")
				}
			}, func() {
				t.Error("customer balked unexpectedly")
			})

			synctest.Wait()

			if state.Load()&customerCutBit == 0 {
				t.Fatal("customer did not begin haircut")
			}
			if state.Load()&barberCutBit == 0 {
				t.Fatal("barber did not begin haircut")
			}

			releaseCutOnce.Do(func() { close(releaseCut) })
			synctest.Wait()

			select {
			case err := <-customerDone:
				if err != nil {
					t.Fatalf("customer returned %v", err)
				}
			default:
				t.Fatal("customer did not complete")
			}

			cancel()
			synctest.Wait()
			requireClosed(t, barberDone, "barber did not exit after cancellation")

			if state.Load()&exitBit == 0 {
				t.Fatal("customer did not exit")
			}
		})
	})

	t.Run("ShopCapacityBalks", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			shop := newImpl()

			var entered atomic.Int32
			var balked atomic.Int32
			done := make([]chan error, 21)

			for i := range 20 {
				done[i] = startCustomer(ctx, shop, func() {
					entered.Add(1)
				}, func() {}, func() {}, func() {}, func() {}, func() {
					t.Error("customer balked before the shop was full")
				})
			}

			synctest.Wait()

			done[20] = startCustomer(ctx, shop, func() {
				t.Error("customer entered even though the shop was full")
			}, func() {}, func() {}, func() {}, func() {}, func() {
				balked.Add(1)
			})

			synctest.Wait()

			if entered.Load() != 20 {
				t.Fatalf("got %d admitted customers, want 20", entered.Load())
			}
			if balked.Load() != 1 {
				t.Fatalf("got %d balked customers, want 1", balked.Load())
			}

			select {
			case err := <-done[20]:
				if err != nil {
					t.Fatalf("balking customer returned %v", err)
				}
			default:
				t.Fatal("balking customer did not complete")
			}

			cancel()
			synctest.Wait()

			for i := range 20 {
				select {
				case err := <-done[i]:
					if err != nil && err != context.Canceled {
						t.Fatalf("customer %d returned %v", i, err)
					}
				default:
					t.Fatalf("customer %d did not exit after cancellation", i)
				}
			}
		})
	})

	t.Run("StandingFIFOToSofa", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			shop := newImpl()

			var mu sync.Mutex
			var sofaOrder []int
			done := make([]chan error, 6)

			for id := 1; id <= 6; id++ {
				id := id
				done[id-1] = startCustomer(ctx, shop, func() {}, func() {
					mu.Lock()
					sofaOrder = append(sofaOrder, id)
					mu.Unlock()
				}, func() {}, func() {}, func() {}, func() {
					t.Errorf("customer %d balked unexpectedly", id)
				})
				synctest.Wait()
			}

			releaseCuts := make(chan struct{})
			var releaseCutsOnce sync.Once
			defer releaseCutsOnce.Do(func() { close(releaseCuts) })
			barberDone1 := startBarber(ctx, shop, func() { <-releaseCuts }, func() {})
			barberDone2 := startBarber(ctx, shop, func() { <-releaseCuts }, func() {})

			synctest.Wait()

			mu.Lock()
			got := append([]int(nil), sofaOrder...)
			mu.Unlock()

			want := []int{1, 2, 3, 4, 5, 6}
			if len(got) != len(want) {
				t.Fatalf("got sofa order %v, want %v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("got sofa order %v, want %v", got, want)
				}
			}

			releaseCutsOnce.Do(func() { close(releaseCuts) })
			synctest.Wait()

			cancel()
			synctest.Wait()
			requireClosed(t, barberDone1, "first barber did not exit after cancellation")
			requireClosed(t, barberDone2, "second barber did not exit after cancellation")
			for i, ch := range done {
				select {
				case err := <-ch:
					if err != nil && err != context.Canceled {
						t.Fatalf("customer %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("customer %d did not complete", i+1)
				}
			}
		})
	})

	t.Run("SofaFIFOToChair", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			shop := newImpl()
			var mu sync.Mutex
			var haircutOrder []int
			release := []chan struct{}{make(chan struct{}), make(chan struct{}), make(chan struct{})}
			var releaseOnce [3]sync.Once
			for i, ch := range release {
				i, ch := i, ch
				defer releaseOnce[i].Do(func() { close(ch) })
			}

			done := make([]chan error, 3)

			for id := 1; id <= 3; id++ {
				id := id
				done[id-1] = startCustomer(ctx, shop, func() {}, func() {}, func() {
					mu.Lock()
					haircutOrder = append(haircutOrder, id)
					mu.Unlock()
					<-release[id-1]
				}, func() {}, func() {}, func() {
					t.Errorf("customer %d balked unexpectedly", id)
				})
				synctest.Wait()
			}

			barberDone := startBarber(ctx, shop, func() {}, func() {})

			synctest.Wait()

			mu.Lock()
			if len(haircutOrder) != 1 || haircutOrder[0] != 1 {
				t.Fatalf("got initial haircut order %v, want [1]", haircutOrder)
			}
			mu.Unlock()

			releaseOnce[0].Do(func() { close(release[0]) })
			synctest.Wait()

			mu.Lock()
			if len(haircutOrder) != 2 || haircutOrder[1] != 2 {
				t.Fatalf("got haircut order %v after first release, want [1 2]", haircutOrder)
			}
			mu.Unlock()

			releaseOnce[1].Do(func() { close(release[1]) })
			synctest.Wait()

			mu.Lock()
			if len(haircutOrder) != 3 || haircutOrder[2] != 3 {
				t.Fatalf("got haircut order %v after second release, want [1 2 3]", haircutOrder)
			}
			mu.Unlock()

			releaseOnce[2].Do(func() { close(release[2]) })
			synctest.Wait()

			for i, ch := range done {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("customer %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("customer %d did not complete", i+1)
				}
			}

			cancel()
			synctest.Wait()
			requireClosed(t, barberDone, "barber did not exit after cancellation")
		})
	})

	t.Run("ThreeConcurrentHaircuts", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			shop := newImpl()

			var activeCuts atomic.Int32
			var activeHaircuts atomic.Int32
			var totalCuts atomic.Int32
			var totalHaircuts atomic.Int32

			release := make(chan struct{})
			var releaseOnce sync.Once
			defer releaseOnce.Do(func() { close(release) })
			barberDone := []chan struct{}{
				startBarber(ctx, shop, func() {
					activeCuts.Add(1)
					totalCuts.Add(1)
					<-release
					activeCuts.Add(-1)
				}, func() {}),
				startBarber(ctx, shop, func() {
					activeCuts.Add(1)
					totalCuts.Add(1)
					<-release
					activeCuts.Add(-1)
				}, func() {}),
				startBarber(ctx, shop, func() {
					activeCuts.Add(1)
					totalCuts.Add(1)
					<-release
					activeCuts.Add(-1)
				}, func() {}),
			}

			done := make([]chan error, 4)
			for i := range 4 {
				done[i] = startCustomer(ctx, shop, func() {}, func() {}, func() {
					activeHaircuts.Add(1)
					totalHaircuts.Add(1)
					<-release
					activeHaircuts.Add(-1)
				}, func() {}, func() {}, func() {
					t.Errorf("customer %d balked unexpectedly", i+1)
				})
			}

			synctest.Wait()

			if activeCuts.Load() != 3 {
				t.Fatalf("got %d concurrent barber cuts, want 3", activeCuts.Load())
			}
			if activeHaircuts.Load() != 3 {
				t.Fatalf("got %d concurrent customer haircuts, want 3", activeHaircuts.Load())
			}
			if totalCuts.Load() != 3 || totalHaircuts.Load() != 3 {
				t.Fatalf("got totalCuts=%d totalHaircuts=%d before release, want 3 and 3", totalCuts.Load(), totalHaircuts.Load())
			}

			releaseOnce.Do(func() { close(release) })
			synctest.Wait()

			if totalCuts.Load() != 4 || totalHaircuts.Load() != 4 {
				t.Fatalf("got totalCuts=%d totalHaircuts=%d after release, want 4 and 4", totalCuts.Load(), totalHaircuts.Load())
			}

			for i, ch := range done {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("customer %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("customer %d did not complete", i+1)
				}
			}

			cancel()
			synctest.Wait()
			for i, ch := range barberDone {
				requireClosed(t, ch, "barber did not exit after cancellation")
				_ = i
			}
		})
	})

	t.Run("SerializedPaymentAndExitAfterReceipt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			shop := newImpl()

			var activeAccept atomic.Int32
			var maxAccept atomic.Int32
			var exits atomic.Int32

			firstAcceptStarted := make(chan struct{})
			var once sync.Once
			releaseAccept := make(chan struct{})
			var releaseAcceptOnce sync.Once
			defer releaseAcceptOnce.Do(func() { close(releaseAccept) })

			barberDone := []chan struct{}{
				startBarber(ctx, shop, func() {}, func() {
					current := activeAccept.Add(1)
					for {
						max := maxAccept.Load()
						if current <= max || maxAccept.CompareAndSwap(max, current) {
							break
						}
					}
					once.Do(func() { close(firstAcceptStarted) })
					<-releaseAccept
					activeAccept.Add(-1)
				}),
				startBarber(ctx, shop, func() {}, func() {
					current := activeAccept.Add(1)
					for {
						max := maxAccept.Load()
						if current <= max || maxAccept.CompareAndSwap(max, current) {
							break
						}
					}
					once.Do(func() { close(firstAcceptStarted) })
					<-releaseAccept
					activeAccept.Add(-1)
				}),
			}

			done := []chan error{
				startCustomer(ctx, shop, func() {}, func() {}, func() {}, func() {}, func() {
					exits.Add(1)
				}, func() { t.Error("customer 1 balked unexpectedly") }),
				startCustomer(ctx, shop, func() {}, func() {}, func() {}, func() {}, func() {
					exits.Add(1)
				}, func() { t.Error("customer 2 balked unexpectedly") }),
			}

			<-firstAcceptStarted

			if exits.Load() != 0 {
				t.Fatalf("got %d exiting customers before receipt was accepted, want 0", exits.Load())
			}
			if activeAccept.Load() != 1 {
				t.Fatalf("got %d concurrent acceptPayment callbacks, want 1", activeAccept.Load())
			}

			releaseAcceptOnce.Do(func() { close(releaseAccept) })
			synctest.Wait()

			if maxAccept.Load() != 1 {
				t.Fatalf("got max %d concurrent acceptPayment callbacks, want 1", maxAccept.Load())
			}
			if exits.Load() != 2 {
				t.Fatalf("got %d exited customers, want 2", exits.Load())
			}

			for i, ch := range done {
				select {
				case err := <-ch:
					if err != nil {
						t.Fatalf("customer %d returned %v", i+1, err)
					}
				default:
					t.Fatalf("customer %d did not complete", i+1)
				}
			}

			cancel()
			synctest.Wait()
			for _, ch := range barberDone {
				requireClosed(t, ch, "barber did not exit after cancellation")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() Barbershop) {
	b.Helper()
	b.ReportAllocs()

	b.Run("ThreeBarbers", func(b *testing.B) {
		ctx := b.Context()
		shop := newImpl()

		for range 3 {
			go func() {
				_ = shop.Barber(ctx, func() {}, func() {})
			}()
		}

		for b.Loop() {
			done := make(chan error, 1)
			go func() {
				done <- shop.Customer(ctx, func() {}, func() {}, func() {}, func() {}, func() {}, func() {})
			}()

			if err := <-done; err != nil {
				b.Fatalf("customer returned %v", err)
			}
		}
	})
}
