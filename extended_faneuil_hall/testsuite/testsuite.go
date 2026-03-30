// Package testsuite contains reusable behavioral tests for extended Faneuil Hall implementations.
package testsuite

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type ExtendedFaneuilHall interface {
	Immigrant(ctx context.Context, enter, checkIn, sitDown, swear, getCertificate, leave func()) error
	Judge(ctx context.Context, enter, confirm, leave func()) error
	Spectator(ctx context.Context, enter, spectate, leave func()) error
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

func requireSuccess(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("%s: %v", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func requireCanceled(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: got %v, want context.Canceled", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func Run(t *testing.T, newImpl func() ExtendedFaneuilHall) {
	t.Helper()

	t.Run("JudgeWaitsForCheckIn", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fh := newImpl()
			ctx := t.Context()

			releaseCheckIn := make(chan struct{})
			immigrantDone := make(chan error, 1)
			judgeDone := make(chan error, 1)
			judgeEntered := make(chan struct{})
			judgeConfirmed := make(chan struct{})

			go func() {
				immigrantDone <- fh.Immigrant(ctx, func() {}, func() {
					<-releaseCheckIn
				}, func() {}, func() {}, func() {}, func() {})
			}()

			synctest.Wait()

			go func() {
				judgeDone <- fh.Judge(ctx, func() {
					close(judgeEntered)
				}, func() {
					close(judgeConfirmed)
				}, func() {})
			}()

			synctest.Wait()
			requireClosed(t, judgeEntered, "judge did not enter")
			requireOpen(t, judgeConfirmed, "judge confirmed before the immigrant checked in")

			close(releaseCheckIn)
			synctest.Wait()

			requireClosed(t, judgeConfirmed, "judge did not confirm after all immigrants checked in")
			requireSuccess(t, judgeDone, "judge did not complete")
			requireSuccess(t, immigrantDone, "immigrant did not complete")
		})
	})

	t.Run("SpectatorMayLeaveWhileJudgePresent", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fh := newImpl()
			ctx := t.Context()

			spectatorEntered := make(chan struct{})
			releaseSpectate := make(chan struct{})
			spectatorLeft := make(chan struct{})
			spectatorDone := make(chan error, 1)
			judgeEntered := make(chan struct{})
			releaseJudge := make(chan struct{})
			judgeDone := make(chan error, 1)

			go func() {
				spectatorDone <- fh.Spectator(ctx, func() {
					close(spectatorEntered)
				}, func() {
					<-releaseSpectate
				}, func() {
					close(spectatorLeft)
				})
			}()

			synctest.Wait()
			requireClosed(t, spectatorEntered, "spectator did not enter")

			go func() {
				judgeDone <- fh.Judge(ctx, func() {
					close(judgeEntered)
				}, func() {}, func() {
					<-releaseJudge
				})
			}()

			synctest.Wait()
			requireClosed(t, judgeEntered, "judge did not enter")

			close(releaseSpectate)
			synctest.Wait()
			requireClosed(t, spectatorLeft, "spectator did not leave while the judge was present")

			close(releaseJudge)
			synctest.Wait()

			requireSuccess(t, spectatorDone, "spectator did not complete")
			requireSuccess(t, judgeDone, "judge did not complete")
		})
	})

	t.Run("ImmigrantCannotLeaveBeforeJudgeLeaves", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fh := newImpl()
			ctx := t.Context()

			immigrantLeft := make(chan struct{})
			immigrantDone := make(chan error, 1)
			confirmed := atomic.Bool{}
			releaseJudge := make(chan struct{})
			judgeDone := make(chan error, 1)

			go func() {
				immigrantDone <- fh.Immigrant(ctx, func() {}, func() {}, func() {}, func() {}, func() {
					if confirmed.Load() {
						return
					}
					t.Error("immigrant got certificate before confirmation")
				}, func() {
					close(immigrantLeft)
				})
			}()

			synctest.Wait()

			go func() {
				judgeDone <- fh.Judge(ctx, func() {}, func() {
					confirmed.Store(true)
				}, func() {
					<-releaseJudge
				})
			}()

			synctest.Wait()
			requireOpen(t, immigrantLeft, "immigrant left while the judge was still in the building")

			close(releaseJudge)
			synctest.Wait()

			requireClosed(t, immigrantLeft, "immigrant did not leave after the judge left")
			requireSuccess(t, judgeDone, "judge did not complete")
			requireSuccess(t, immigrantDone, "immigrant did not complete")
		})
	})

	t.Run("EntryBlockedWhileJudgePresent", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fh := newImpl()
			ctx := t.Context()

			immigrantDone := make(chan error, 1)
			releaseJudge := make(chan struct{})
			judgeDone := make(chan error, 1)
			spectatorEntered := make(chan struct{})
			spectatorDone := make(chan error, 1)

			go func() {
				immigrantDone <- fh.Immigrant(ctx, func() {}, func() {}, func() {}, func() {}, func() {}, func() {})
			}()

			synctest.Wait()

			go func() {
				judgeDone <- fh.Judge(ctx, func() {}, func() {}, func() {
					<-releaseJudge
				})
			}()

			synctest.Wait()

			go func() {
				spectatorDone <- fh.Spectator(ctx, func() {
					close(spectatorEntered)
				}, func() {}, func() {})
			}()

			synctest.Wait()
			requireOpen(t, spectatorEntered, "spectator entered while the judge was present")

			close(releaseJudge)
			synctest.Wait()

			requireClosed(t, spectatorEntered, "spectator did not enter after the judge left")
			requireSuccess(t, judgeDone, "judge did not complete")
			requireSuccess(t, immigrantDone, "immigrant did not complete")
			requireSuccess(t, spectatorDone, "spectator did not complete")
		})
	})

	t.Run("NextJudgeWaitsForSwornImmigrantsToLeave", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fh := newImpl()
			ctx := t.Context()

			releaseImmigrantLeave := make(chan struct{})
			immigrantLeft := make(chan struct{})
			immigrantDone := make(chan error, 1)
			firstJudgeDone := make(chan error, 1)
			secondJudgeEntered := make(chan struct{})
			secondJudgeDone := make(chan error, 1)

			go func() {
				immigrantDone <- fh.Immigrant(ctx, func() {}, func() {}, func() {}, func() {}, func() {
					<-releaseImmigrantLeave
				}, func() {
					close(immigrantLeft)
				})
			}()

			synctest.Wait()

			go func() {
				firstJudgeDone <- fh.Judge(ctx, func() {}, func() {}, func() {})
			}()

			synctest.Wait()
			requireSuccess(t, firstJudgeDone, "first judge did not complete")
			requireOpen(t, immigrantLeft, "immigrant left before being released to do so")

			go func() {
				secondJudgeDone <- fh.Judge(ctx, func() {
					close(secondJudgeEntered)
				}, func() {}, func() {})
			}()

			synctest.Wait()
			requireOpen(t, secondJudgeEntered, "second judge entered before the sworn immigrant left")

			close(releaseImmigrantLeave)
			synctest.Wait()

			requireClosed(t, immigrantLeft, "immigrant did not leave after being released")
			requireClosed(t, secondJudgeEntered, "second judge did not enter after the sworn immigrant left")
			requireSuccess(t, immigrantDone, "immigrant did not complete")
			requireSuccess(t, secondJudgeDone, "second judge did not complete")
		})
	})

	t.Run("CancelWaitingEntrant", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			run  func(fh ExtendedFaneuilHall, ctx context.Context) <-chan error
		}{
			{
				name: "Immigrant",
				run: func(fh ExtendedFaneuilHall, ctx context.Context) <-chan error {
					done := make(chan error, 1)
					go func() {
						done <- fh.Immigrant(ctx, func() {}, func() {}, func() {}, func() {}, func() {}, func() {})
					}()
					return done
				},
			},
			{
				name: "Spectator",
				run: func(fh ExtendedFaneuilHall, ctx context.Context) <-chan error {
					done := make(chan error, 1)
					go func() {
						done <- fh.Spectator(ctx, func() {}, func() {}, func() {})
					}()
					return done
				},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					fh := newImpl()
					ctx := t.Context()

					immigrantDone := make(chan error, 1)
					releaseJudge := make(chan struct{})
					go func() {
						immigrantDone <- fh.Immigrant(ctx, func() {}, func() {}, func() {}, func() {}, func() {}, func() {})
					}()

					synctest.Wait()

					go fh.Judge(ctx, func() {}, func() {}, func() {
						<-releaseJudge
					})

					synctest.Wait()

					waitCtx, cancel := context.WithCancel(ctx)
					defer cancel()
					waiterDone := tc.run(fh, waitCtx)

					synctest.Wait()
					cancel()
					synctest.Wait()

					requireCanceled(t, waiterDone, tc.name+" did not exit after cancellation while waiting to enter")

					close(releaseJudge)
					synctest.Wait()
					requireSuccess(t, immigrantDone, "immigrant did not complete")
				})
			})
		}
	})

	t.Run("CancelWaitingJudge", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fh := newImpl()
			ctx := t.Context()

			releaseImmigrantLeave := make(chan struct{})
			immigrantDone := make(chan error, 1)
			firstJudgeDone := make(chan error, 1)
			secondJudgeDone := make(chan error, 1)

			go func() {
				immigrantDone <- fh.Immigrant(ctx, func() {}, func() {}, func() {}, func() {}, func() {
					<-releaseImmigrantLeave
				}, func() {})
			}()

			synctest.Wait()

			go func() {
				firstJudgeDone <- fh.Judge(ctx, func() {}, func() {}, func() {})
			}()

			synctest.Wait()
			requireSuccess(t, firstJudgeDone, "first judge did not complete")

			judgeCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			go func() {
				secondJudgeDone <- fh.Judge(judgeCtx, func() {}, func() {}, func() {})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			requireCanceled(t, secondJudgeDone, "second judge did not exit after cancellation while waiting")

			close(releaseImmigrantLeave)
			synctest.Wait()
			requireSuccess(t, immigrantDone, "immigrant did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() ExtendedFaneuilHall) {
	b.Helper()
	b.ReportAllocs()

	b.Run("OneImmigrantOneJudge", func(b *testing.B) {
		for b.Loop() {
			fh := newImpl()
			ctx := b.Context()
			entered := make(chan struct{})

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				if err := fh.Immigrant(ctx, func() {
					close(entered)
				}, func() {}, func() {}, func() {}, func() {}, func() {}); err != nil {
					b.Errorf("immigrant returned %v", err)
				}
			}()

			go func() {
				defer wg.Done()
				<-entered
				if err := fh.Judge(ctx, func() {}, func() {}, func() {}); err != nil {
					b.Errorf("judge returned %v", err)
				}
			}()

			wg.Wait()
		}
	})

	b.Run("ImmigrantSpectatorJudge", func(b *testing.B) {
		for b.Loop() {
			fh := newImpl()
			ctx := b.Context()
			entered := make(chan struct{})

			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				if err := fh.Immigrant(ctx, func() {
					close(entered)
				}, func() {}, func() {}, func() {}, func() {}, func() {}); err != nil {
					b.Errorf("immigrant returned %v", err)
				}
			}()

			go func() {
				defer wg.Done()
				if err := fh.Spectator(ctx, func() {}, func() {}, func() {}); err != nil {
					b.Errorf("spectator returned %v", err)
				}
			}()

			go func() {
				defer wg.Done()
				<-entered
				if err := fh.Judge(ctx, func() {}, func() {}, func() {}); err != nil {
					b.Errorf("judge returned %v", err)
				}
			}()

			wg.Wait()
		}
	})
}
