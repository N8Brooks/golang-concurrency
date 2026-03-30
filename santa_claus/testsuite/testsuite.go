// Package testsuite contains reusable behavioral tests for Santa Claus implementations.
package testsuite

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const numReindeer = 9

type SantaClaus interface {
	Santa(ctx context.Context, prepareSleigh, helpElves func())
	Reindeer(ctx context.Context, getHitched func())
	Elf(ctx context.Context, getHelp func())
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

func Run(t *testing.T, newImpl func() SantaClaus) {
	t.Helper()

	t.Run("ReindeerBatch", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sc := newImpl()

			doneSanta := make(chan struct{})
			var prepareCalls atomic.Int64
			var prepared atomic.Bool
			var hitched atomic.Int64

			go func() {
				sc.Santa(ctx, func() {
					if prepareCalls.Add(1) != 1 {
						t.Error("prepareSleigh called more than once for one reindeer batch")
					}
					prepared.Store(true)
				}, func() {
					t.Error("helpElves called without three elves waiting")
				})
				close(doneSanta)
			}()

			doneReindeer := make([]chan struct{}, numReindeer)
			for i := range numReindeer - 1 {
				doneReindeer[i] = make(chan struct{})
				go func(done chan struct{}) {
					sc.Reindeer(ctx, func() {
						if !prepared.Load() {
							t.Error("getHitched ran before prepareSleigh completed")
						}
						hitched.Add(1)
					})
					close(done)
				}(doneReindeer[i])
			}

			synctest.Wait()

			if prepareCalls.Load() != 0 {
				t.Fatalf("prepareSleigh called %d times before the ninth reindeer arrived", prepareCalls.Load())
			}
			if hitched.Load() != 0 {
				t.Fatalf("getHitched ran %d times before the ninth reindeer arrived", hitched.Load())
			}

			doneReindeer[numReindeer-1] = make(chan struct{})
			go func(done chan struct{}) {
				sc.Reindeer(ctx, func() {
					if !prepared.Load() {
						t.Error("getHitched ran before prepareSleigh completed")
					}
					hitched.Add(1)
				})
				close(done)
			}(doneReindeer[numReindeer-1])

			synctest.Wait()

			if prepareCalls.Load() != 1 {
				t.Fatalf("prepareSleigh called %d times, want 1", prepareCalls.Load())
			}
			if hitched.Load() != numReindeer {
				t.Fatalf("getHitched ran %d times, want %d", hitched.Load(), numReindeer)
			}

			for i, done := range doneReindeer {
				requireClosed(t, done, "reindeer "+string(rune('1'+i))+" did not return after getHitched")
			}

			cancel()
			synctest.Wait()
			requireClosed(t, doneSanta, "Santa did not exit after cancellation")
		})
	})

	t.Run("ElfBatching", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sc := newImpl()
			doneSanta := make(chan struct{})
			var helpCalls atomic.Int64
			var helped atomic.Int64
			releaseFirstBatch := make(chan struct{})

			go func() {
				sc.Santa(ctx, func() {
					t.Error("prepareSleigh called without nine reindeer")
				}, func() {
					helpCalls.Add(1)
				})
				close(doneSanta)
			}()

			doneElves := make([]chan struct{}, 6)
			for i := range doneElves {
				doneElves[i] = make(chan struct{})
				go func(done chan struct{}) {
					sc.Elf(ctx, func() {
						if n := helped.Add(1); n <= 3 {
							<-releaseFirstBatch
						}
					})
					close(done)
				}(doneElves[i])
			}

			synctest.Wait()

			if helpCalls.Load() != 1 {
				t.Fatalf("helpElves called %d times while the first batch was in progress, want 1", helpCalls.Load())
			}
			if helped.Load() != 3 {
				t.Fatalf("getHelp ran %d times while the first batch was in progress, want 3", helped.Load())
			}

			close(releaseFirstBatch)
			synctest.Wait()

			if helpCalls.Load() != 2 {
				t.Fatalf("helpElves called %d times after two elf batches, want 2", helpCalls.Load())
			}
			if helped.Load() != 6 {
				t.Fatalf("getHelp ran %d times after two elf batches, want 6", helped.Load())
			}

			for _, done := range doneElves {
				requireClosed(t, done, "an elf did not return after getHelp")
			}

			cancel()
			synctest.Wait()
			requireClosed(t, doneSanta, "Santa did not exit after cancellation")
		})
	})

	t.Run("ReindeerPriority", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sc := newImpl()
			releasePrepare := make(chan struct{})
			releaseElves := make(chan struct{})
			prepareStarted := make(chan struct{})
			helpStarted := make(chan struct{})
			doneSanta := make(chan struct{})

			for range 3 {
				go sc.Elf(ctx, func() {
					<-releaseElves
				})
			}
			for range numReindeer {
				go sc.Reindeer(ctx, func() {})
			}

			synctest.Wait()

			go func() {
				sc.Santa(ctx, func() {
					close(prepareStarted)
					<-releasePrepare
				}, func() {
					close(helpStarted)
				})
				close(doneSanta)
			}()

			synctest.Wait()

			requireClosed(t, prepareStarted, "Santa did not prioritize the reindeer when both groups were ready")
			requireOpen(t, helpStarted, "Santa helped elves before preparing the sleigh")

			close(releasePrepare)
			synctest.Wait()
			requireClosed(t, helpStarted, "Santa did not help the waiting elves after the reindeer batch")

			close(releaseElves)
			cancel()
			synctest.Wait()
			requireClosed(t, doneSanta, "Santa did not exit after cancellation")
		})
	})

	t.Run("CancelIdleSanta", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sc := newImpl()
			done := make(chan struct{})

			go func() {
				sc.Santa(ctx, func() {}, func() {})
				close(done)
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			requireClosed(t, done, "Santa did not exit after cancellation while idle")
		})
	})

	t.Run("CancelWaitingElf", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sc := newImpl()
			go sc.Santa(ctx, func() {}, func() {})

			elfCtx, cancelElf := context.WithCancel(ctx)
			defer cancelElf()

			var helped atomic.Int64
			done := make(chan struct{})

			go func() {
				sc.Elf(elfCtx, func() {
					helped.Add(1)
				})
				close(done)
			}()

			synctest.Wait()
			cancelElf()
			synctest.Wait()

			requireClosed(t, done, "elf did not exit after cancellation while waiting")
			if helped.Load() != 0 {
				t.Fatalf("getHelp ran %d times for a canceled waiting elf, want 0", helped.Load())
			}
		})
	})

	t.Run("CancelWaitingReindeer", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sc := newImpl()
			go sc.Santa(ctx, func() {}, func() {})

			reindeerCtx, cancelReindeer := context.WithCancel(ctx)
			defer cancelReindeer()

			var hitched atomic.Int64
			done := make(chan struct{})

			go func() {
				sc.Reindeer(reindeerCtx, func() {
					hitched.Add(1)
				})
				close(done)
			}()

			synctest.Wait()
			cancelReindeer()
			synctest.Wait()

			requireClosed(t, done, "reindeer did not exit after cancellation while waiting")
			if hitched.Load() != 0 {
				t.Fatalf("getHitched ran %d times for a canceled waiting reindeer, want 0", hitched.Load())
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() SantaClaus) {
	b.Helper()
	b.ReportAllocs()

	b.Run("ElfBatch", func(b *testing.B) {
		for b.Loop() {
			ctx, cancel := context.WithCancel(b.Context())
			sc := newImpl()

			var wg sync.WaitGroup
			wg.Add(4)

			go func() {
				defer wg.Done()
				sc.Santa(ctx, func() {}, func() {
					cancel()
				})
			}()

			for range 3 {
				go func() {
					defer wg.Done()
					sc.Elf(ctx, func() {})
				}()
			}

			wg.Wait()
		}
	})

	b.Run("ReindeerBatch", func(b *testing.B) {
		for b.Loop() {
			ctx, cancel := context.WithCancel(b.Context())
			sc := newImpl()

			var wg sync.WaitGroup
			var hitched atomic.Int64
			wg.Add(numReindeer + 1)

			go func() {
				defer wg.Done()
				sc.Santa(ctx, func() {}, func() {})
			}()

			for range numReindeer {
				go func() {
					defer wg.Done()
					sc.Reindeer(ctx, func() {
						if hitched.Add(1) == numReindeer {
							cancel()
						}
					})
				}()
			}

			wg.Wait()
		}
	})
}
