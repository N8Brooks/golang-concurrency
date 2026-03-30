// Package testsuite contains reusable behavioral tests for generalized smokers implementations.
package testsuite

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"math/rand/v2"
)

type Agent interface {
	Tobacco() chan struct{}
	Paper() chan struct{}
	Match() chan struct{}
}

type Smokers interface {
	Run(ctx context.Context)
	SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func())
	SmokerWithPaper(ctx context.Context, makeCigarette, smoke func())
	SmokerWithMatch(ctx context.Context, makeCigarette, smoke func())
}

type agent struct {
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
}

func newAgent(buffer int) *agent {
	return &agent{
		tobacco: make(chan struct{}, buffer),
		paper:   make(chan struct{}, buffer),
		match:   make(chan struct{}, buffer),
	}
}

func (a *agent) Tobacco() chan struct{} {
	return a.tobacco
}

func (a *agent) Paper() chan struct{} {
	return a.paper
}

func (a *agent) Match() chan struct{} {
	return a.match
}

func startLoops(ctx context.Context, s Smokers, tSmoke, pSmoke, mSmoke func()) [4]chan struct{} {
	doneRun := make(chan struct{})
	doneTobacco := make(chan struct{})
	donePaper := make(chan struct{})
	doneMatch := make(chan struct{})

	go func() {
		s.Run(ctx)
		close(doneRun)
	}()
	go func() {
		s.SmokerWithTobacco(ctx, func() {}, tSmoke)
		close(doneTobacco)
	}()
	go func() {
		s.SmokerWithPaper(ctx, func() {}, pSmoke)
		close(donePaper)
	}()
	go func() {
		s.SmokerWithMatch(ctx, func() {}, mSmoke)
		close(doneMatch)
	}()

	return [4]chan struct{}{doneRun, doneTobacco, donePaper, doneMatch}
}

func requireClosed(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	default:
		t.Fatal(msg)
	}
}

func Run(t *testing.T, newImpl func(a Agent) Smokers) {
	t.Helper()

	t.Run("TargetedPairs", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			a := newAgent(8)
			cs := newImpl(a)

			var tobaccoSmokes, paperSmokes, matchSmokes atomic.Int64
			done := startLoops(ctx, cs, func() {
				tobaccoSmokes.Add(1)
			}, func() {
				paperSmokes.Add(1)
			}, func() {
				matchSmokes.Add(1)
			})

			a.paper <- struct{}{}
			a.match <- struct{}{}
			synctest.Wait()

			if tobaccoSmokes.Load() != 1 || paperSmokes.Load() != 0 || matchSmokes.Load() != 0 {
				t.Fatalf("paper+match should wake only smoker with tobacco, got tobacco=%d paper=%d match=%d", tobaccoSmokes.Load(), paperSmokes.Load(), matchSmokes.Load())
			}

			a.tobacco <- struct{}{}
			a.match <- struct{}{}
			synctest.Wait()

			if tobaccoSmokes.Load() != 1 || paperSmokes.Load() != 1 || matchSmokes.Load() != 0 {
				t.Fatalf("tobacco+match should wake only smoker with paper, got tobacco=%d paper=%d match=%d", tobaccoSmokes.Load(), paperSmokes.Load(), matchSmokes.Load())
			}

			a.tobacco <- struct{}{}
			a.paper <- struct{}{}
			synctest.Wait()

			if tobaccoSmokes.Load() != 1 || paperSmokes.Load() != 1 || matchSmokes.Load() != 1 {
				t.Fatalf("tobacco+paper should wake only smoker with match, got tobacco=%d paper=%d match=%d", tobaccoSmokes.Load(), paperSmokes.Load(), matchSmokes.Load())
			}

			cancel()
			synctest.Wait()
			requireClosed(t, done[0], "Run did not exit after cancellation")
			requireClosed(t, done[1], "SmokerWithTobacco did not exit after cancellation")
			requireClosed(t, done[2], "SmokerWithPaper did not exit after cancellation")
			requireClosed(t, done[3], "SmokerWithMatch did not exit after cancellation")
		})
	})

	t.Run("AccumulatedSupplies", func(t *testing.T) {
		const numIterations = 100

		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			a := newAgent(numIterations * 2)
			cs := newImpl(a)

			var tobaccos, papers, matches atomic.Int64
			var tobaccoInFlight, paperInFlight, matchInFlight atomic.Int64
			var tobaccoSmoked, paperSmoked, matchSmoked atomic.Int64

			doneRun := make(chan struct{})
			go func() {
				cs.Run(ctx)
				close(doneRun)
			}()

			doneTobacco := make(chan struct{})
			go func() {
				cs.SmokerWithTobacco(ctx, func() {
					if papers.Add(-1) < 0 || matches.Add(-1) < 0 {
						t.Error("smoker with tobacco could not make a cigarette")
					}
					tobaccoInFlight.Add(1)
				}, func() {
					if tobaccoInFlight.Add(-1) < 0 {
						t.Error("smoker with tobacco could not smoke a cigarette")
					}
					tobaccoSmoked.Add(1)
				})
				close(doneTobacco)
			}()

			donePaper := make(chan struct{})
			go func() {
				cs.SmokerWithPaper(ctx, func() {
					if tobaccos.Add(-1) < 0 || matches.Add(-1) < 0 {
						t.Error("smoker with paper could not make a cigarette")
					}
					paperInFlight.Add(1)
				}, func() {
					if paperInFlight.Add(-1) < 0 {
						t.Error("smoker with paper could not smoke a cigarette")
					}
					paperSmoked.Add(1)
				})
				close(donePaper)
			}()

			doneMatch := make(chan struct{})
			go func() {
				cs.SmokerWithMatch(ctx, func() {
					if tobaccos.Add(-1) < 0 || papers.Add(-1) < 0 {
						t.Error("smoker with match could not make a cigarette")
					}
					matchInFlight.Add(1)
				}, func() {
					if matchInFlight.Add(-1) < 0 {
						t.Error("smoker with match could not smoke a cigarette")
					}
					matchSmoked.Add(1)
				})
				close(doneMatch)
			}()

			rng := rand.New(rand.NewPCG(1, 2))
			var suppliedTobacco, suppliedPaper, suppliedMatch int64

			for range numIterations {
				switch rng.IntN(3) {
				case 0:
					papers.Add(1)
					matches.Add(1)
					suppliedPaper++
					suppliedMatch++
					a.paper <- struct{}{}
					a.match <- struct{}{}
				case 1:
					tobaccos.Add(1)
					matches.Add(1)
					suppliedTobacco++
					suppliedMatch++
					a.tobacco <- struct{}{}
					a.match <- struct{}{}
				case 2:
					tobaccos.Add(1)
					papers.Add(1)
					suppliedTobacco++
					suppliedPaper++
					a.tobacco <- struct{}{}
					a.paper <- struct{}{}
				}
			}

			synctest.Wait()

			leftTobacco := tobaccos.Load()
			leftPaper := papers.Load()
			leftMatch := matches.Load()

			if (leftTobacco > 0 && leftPaper > 0) || (leftTobacco > 0 && leftMatch > 0) || (leftPaper > 0 && leftMatch > 0) {
				t.Fatalf("expected no complete pair to remain, got tobacco=%d paper=%d match=%d", leftTobacco, leftPaper, leftMatch)
			}

			if tobaccoInFlight.Load() != 0 || paperInFlight.Load() != 0 || matchInFlight.Load() != 0 {
				t.Fatalf("expected no in-flight cigarettes, got tobacco=%d paper=%d match=%d", tobaccoInFlight.Load(), paperInFlight.Load(), matchInFlight.Load())
			}

			actualTotal := tobaccoSmoked.Load() + paperSmoked.Load() + matchSmoked.Load()
			expectedTotal := (suppliedTobacco + suppliedPaper + suppliedMatch - leftTobacco - leftPaper - leftMatch) / 2
			if actualTotal != expectedTotal {
				t.Fatalf("got %d smoked cigarettes, want %d", actualTotal, expectedTotal)
			}

			cancel()
			synctest.Wait()
			requireClosed(t, doneRun, "Run did not exit after cancellation")
			requireClosed(t, doneTobacco, "SmokerWithTobacco did not exit after cancellation")
			requireClosed(t, donePaper, "SmokerWithPaper did not exit after cancellation")
			requireClosed(t, doneMatch, "SmokerWithMatch did not exit after cancellation")
		})
	})

	t.Run("CancelWhileIdle", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			a := newAgent(1)
			cs := newImpl(a)
			done := startLoops(ctx, cs, func() {}, func() {}, func() {})

			synctest.Wait()
			cancel()
			synctest.Wait()

			requireClosed(t, done[0], "Run did not exit after cancellation")
			requireClosed(t, done[1], "SmokerWithTobacco did not exit after cancellation")
			requireClosed(t, done[2], "SmokerWithPaper did not exit after cancellation")
			requireClosed(t, done[3], "SmokerWithMatch did not exit after cancellation")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(a Agent) Smokers) {
	b.Helper()
	b.ReportAllocs()

	for _, tc := range []struct {
		name       string
		supply     func(a Agent)
		smokeQueue func(tobaccoSmoked, paperSmoked, matchSmoked chan struct{}) chan struct{}
	}{
		{
			name: "TobaccoSmoker",
			supply: func(a Agent) {
				a.Paper() <- struct{}{}
				a.Match() <- struct{}{}
			},
			smokeQueue: func(tobaccoSmoked, paperSmoked, matchSmoked chan struct{}) chan struct{} {
				return tobaccoSmoked
			},
		},
		{
			name: "PaperSmoker",
			supply: func(a Agent) {
				a.Tobacco() <- struct{}{}
				a.Match() <- struct{}{}
			},
			smokeQueue: func(tobaccoSmoked, paperSmoked, matchSmoked chan struct{}) chan struct{} {
				return paperSmoked
			},
		},
		{
			name: "MatchSmoker",
			supply: func(a Agent) {
				a.Tobacco() <- struct{}{}
				a.Paper() <- struct{}{}
			},
			smokeQueue: func(tobaccoSmoked, paperSmoked, matchSmoked chan struct{}) chan struct{} {
				return matchSmoked
			},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			ctx := b.Context()
			a := newAgent(2)
			cs := newImpl(a)

			tobaccoSmoked := make(chan struct{}, 1)
			paperSmoked := make(chan struct{}, 1)
			matchSmoked := make(chan struct{}, 1)

			go cs.Run(ctx)
			go cs.SmokerWithTobacco(ctx, func() {}, func() {
				tobaccoSmoked <- struct{}{}
			})
			go cs.SmokerWithPaper(ctx, func() {}, func() {
				paperSmoked <- struct{}{}
			})
			go cs.SmokerWithMatch(ctx, func() {}, func() {
				matchSmoked <- struct{}{}
			})

			smoked := tc.smokeQueue(tobaccoSmoked, paperSmoked, matchSmoked)
			for b.Loop() {
				tc.supply(a)
				<-smoked
			}
		})
	}
}
