// Package testsuite contains reusable behavioral tests for exclusive queue implementations.
package testsuite

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const (
	leaderArrivedBit uint32 = 1 << iota
	followerArrivedBit
	leaderDancedBit
	followerDancedBit
)

type ExclusiveQueue interface {
	Leader(ctx context.Context, leader, dance func()) error
	Follower(ctx context.Context, follower, dance func()) error
}

func runOneShot(t *testing.T, q ExclusiveQueue, ctx context.Context, leaderFirst bool) {
	t.Helper()

	var state atomic.Uint32
	leaderDone := make(chan error, 1)
	followerDone := make(chan error, 1)

	runLeader := func() {
		leaderDone <- q.Leader(ctx, func() {
			if state.Or(leaderArrivedBit)&(leaderDancedBit|followerDancedBit) != 0 {
				t.Error("leader arrival ran after dancing began")
			}
		}, func() {
			if state.Or(leaderDancedBit)&(leaderArrivedBit|followerArrivedBit) != (leaderArrivedBit | followerArrivedBit) {
				t.Error("leader danced before both participants arrived")
			}
		})
	}
	runFollower := func() {
		followerDone <- q.Follower(ctx, func() {
			if state.Or(followerArrivedBit)&(leaderDancedBit|followerDancedBit) != 0 {
				t.Error("follower arrival ran after dancing began")
			}
		}, func() {
			if state.Or(followerDancedBit)&(leaderArrivedBit|followerArrivedBit) != (leaderArrivedBit | followerArrivedBit) {
				t.Error("follower danced before both participants arrived")
			}
		})
	}

	if leaderFirst {
		go runLeader()
		synctest.Wait()
		go runFollower()
	} else {
		go runFollower()
		synctest.Wait()
		go runLeader()
	}

	synctest.Wait()

	select {
	case err := <-leaderDone:
		if err != nil {
			t.Fatalf("leader returned %v", err)
		}
	default:
		t.Fatal("leader did not complete")
	}

	select {
	case err := <-followerDone:
		if err != nil {
			t.Fatalf("follower returned %v", err)
		}
	default:
		t.Fatal("follower did not complete")
	}

	if state.Load()&leaderDancedBit == 0 {
		t.Fatal("leader did not dance")
	}
	if state.Load()&followerDancedBit == 0 {
		t.Fatal("follower did not dance")
	}
}

func Run(t *testing.T, newImpl func() ExclusiveQueue) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name        string
			leaderFirst bool
		}{
			{name: "LeaderFirst", leaderFirst: true},
			{name: "FollowerFirst", leaderFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneShot(t, newImpl(), t.Context(), tc.leaderFirst)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		for _, tc := range []struct {
			name        string
			leaderFirst bool
		}{
			{name: "LeaderFirst", leaderFirst: true},
			{name: "FollowerFirst", leaderFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					q := newImpl()
					ctx := t.Context()
					for range 5 {
						runOneShot(t, q, ctx, tc.leaderFirst)
					}
				})
			})
		}
	})

	t.Run("ExclusiveDance", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			q := newImpl()
			ctx := t.Context()

			pair1Entered := atomic.Int32{}
			pair2Entered := atomic.Int32{}
			releasePair1 := make(chan struct{})

			leader1Done := make(chan error, 1)
			follower1Done := make(chan error, 1)
			leader2Done := make(chan error, 1)
			follower2Done := make(chan error, 1)

			go func() {
				leader1Done <- q.Leader(ctx, func() {}, func() {
					pair1Entered.Add(1)
					<-releasePair1
				})
			}()
			go func() {
				follower1Done <- q.Follower(ctx, func() {}, func() {
					pair1Entered.Add(1)
					<-releasePair1
				})
			}()

			synctest.Wait()

			if pair1Entered.Load() != 2 {
				t.Fatalf("got %d first-pair dancers, want 2", pair1Entered.Load())
			}

			go func() {
				leader2Done <- q.Leader(ctx, func() {}, func() {
					pair2Entered.Add(1)
				})
			}()
			go func() {
				follower2Done <- q.Follower(ctx, func() {}, func() {
					pair2Entered.Add(1)
				})
			}()

			synctest.Wait()

			if pair2Entered.Load() != 0 {
				t.Fatalf("got %d second-pair dancers while first pair was still dancing, want 0", pair2Entered.Load())
			}

			close(releasePair1)
			synctest.Wait()

			if pair2Entered.Load() != 2 {
				t.Fatalf("got %d second-pair dancers after first pair finished, want 2", pair2Entered.Load())
			}

			for _, done := range []chan error{leader1Done, follower1Done, leader2Done, follower2Done} {
				select {
				case err := <-done:
					if err != nil {
						t.Fatalf("participant returned %v", err)
					}
				default:
					t.Fatal("participant did not complete")
				}
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		for _, tc := range []struct {
			name          string
			cancelLeader  bool
			arrivedBit    uint32
			dancedBit     uint32
			expectedError string
		}{
			{name: "Leader", cancelLeader: true, arrivedBit: leaderArrivedBit, dancedBit: leaderDancedBit, expectedError: "leader"},
			{name: "Follower", cancelLeader: false, arrivedBit: followerArrivedBit, dancedBit: followerDancedBit, expectedError: "follower"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()

					q := newImpl()
					var state atomic.Uint32
					done := make(chan error, 1)

					if tc.cancelLeader {
						go func() {
							done <- q.Leader(ctx, func() {
								state.Or(leaderArrivedBit)
							}, func() {
								state.Or(leaderDancedBit)
							})
						}()
					} else {
						go func() {
							done <- q.Follower(ctx, func() {
								state.Or(followerArrivedBit)
							}, func() {
								state.Or(followerDancedBit)
							})
						}()
					}

					synctest.Wait()

					cancel()
					synctest.Wait()

					select {
					case err := <-done:
						if !errors.Is(err, context.Canceled) {
							t.Fatalf("%s returned %v, want context.Canceled", tc.expectedError, err)
						}
					default:
						t.Fatalf("%s did not exit after cancellation", tc.expectedError)
					}

					if state.Load()&tc.arrivedBit == 0 {
						t.Fatalf("%s arrival callback did not run", tc.expectedError)
					}
					if state.Load()&tc.dancedBit != 0 {
						t.Fatalf("%s danced despite never being matched", tc.expectedError)
					}
				})
			})
		}
	})

	t.Run("CancelRemovesWaiter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			q := newImpl()

			canceledCtx, cancel := context.WithCancel(t.Context())
			defer cancel()

			canceledLeaderDone := make(chan error, 1)
			go func() {
				canceledLeaderDone <- q.Leader(canceledCtx, func() {}, func() {
					t.Error("canceled leader danced")
				})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			select {
			case err := <-canceledLeaderDone:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("canceled leader returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("canceled leader did not exit")
			}

			followerDone := make(chan error, 1)
			go func() {
				followerDone <- q.Follower(t.Context(), func() {}, func() {})
			}()

			synctest.Wait()

			select {
			case err := <-followerDone:
				t.Fatalf("follower returned early after matching a canceled leader: %v", err)
			default:
			}

			leaderDone := make(chan error, 1)
			go func() {
				leaderDone <- q.Leader(t.Context(), func() {}, func() {})
			}()

			synctest.Wait()

			select {
			case err := <-followerDone:
				if err != nil {
					t.Fatalf("follower returned %v", err)
				}
			default:
				t.Fatal("follower did not complete after a new leader arrived")
			}

			select {
			case err := <-leaderDone:
				if err != nil {
					t.Fatalf("leader returned %v", err)
				}
			default:
				t.Fatal("leader did not complete after pairing with follower")
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() ExclusiveQueue) {
	b.Helper()
	b.ReportAllocs()

	b.Run("LeaderFirst", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			q := newImpl()
			done := make(chan error, 1)

			go func() {
				done <- q.Leader(ctx, func() {}, func() {})
			}()

			if err := q.Follower(ctx, func() {}, func() {}); err != nil {
				b.Fatalf("follower failed: %v", err)
			}
			if err := <-done; err != nil {
				b.Fatalf("leader failed: %v", err)
			}
		}
	})

	b.Run("FollowerFirst", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			q := newImpl()
			done := make(chan error, 1)

			go func() {
				done <- q.Follower(ctx, func() {}, func() {})
			}()

			if err := q.Leader(ctx, func() {}, func() {}); err != nil {
				b.Fatalf("leader failed: %v", err)
			}
			if err := <-done; err != nil {
				b.Fatalf("follower failed: %v", err)
			}
		}
	})
}
