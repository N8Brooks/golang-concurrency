// Package testsuite contains reusable behavioral tests for cigarette smokers implementations.
package testsuite

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type Agent interface {
	Tobacco() chan struct{}
	Paper() chan struct{}
	Match() chan struct{}
	Signal()
}

type CigaretteSmokers interface {
	Run(ctx context.Context)
	SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func())
	SmokerWithPaper(ctx context.Context, makeCigarette, smoke func())
	SmokerWithMatch(ctx context.Context, makeCigarette, smoke func())
}

type smokerID int

const (
	smokerWithTobacco smokerID = iota + 1
	smokerWithPaper
	smokerWithMatch
)

type ingredientPair int

const (
	tobaccoPaper ingredientPair = iota
	paperMatch
	tobaccoMatch
)

func (p ingredientPair) String() string {
	switch p {
	case tobaccoPaper:
		return "TobaccoPaper"
	case paperMatch:
		return "PaperMatch"
	case tobaccoMatch:
		return "TobaccoMatch"
	default:
		return fmt.Sprintf("ingredientPair(%d)", p)
	}
}

func (p ingredientPair) expectedSmoker() smokerID {
	switch p {
	case tobaccoPaper:
		return smokerWithMatch
	case paperMatch:
		return smokerWithTobacco
	case tobaccoMatch:
		return smokerWithPaper
	default:
		panic("unexpected ingredient pair")
	}
}

type agent struct {
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
	agent   chan struct{}
}

func newAgent() *agent {
	a := &agent{
		tobacco: make(chan struct{}, 1),
		paper:   make(chan struct{}, 1),
		match:   make(chan struct{}, 1),
		agent:   make(chan struct{}, 1),
	}
	a.agent <- struct{}{}
	return a
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

func (a *agent) Signal() {
	a.agent <- struct{}{}
}

func (a *agent) beginRound(t *testing.T, pair ingredientPair) {
	t.Helper()

	select {
	case <-a.agent:
	default:
		t.Fatal("agent was not ready for the next round")
	}

	switch pair {
	case tobaccoPaper:
		a.tobacco <- struct{}{}
		a.paper <- struct{}{}
	case paperMatch:
		a.paper <- struct{}{}
		a.match <- struct{}{}
	case tobaccoMatch:
		a.tobacco <- struct{}{}
		a.match <- struct{}{}
	default:
		t.Fatalf("unexpected pair %v", pair)
	}
}

func (a *agent) waitRoundComplete(t *testing.T) {
	t.Helper()

	select {
	case <-a.agent:
	default:
		t.Fatal("agent was not signaled after a complete round")
	}

	a.agent <- struct{}{}
}

type counters struct {
	tobaccoMade   int64
	paperMade     int64
	matchMade     int64
	tobaccoSmoked int64
	paperSmoked   int64
	matchSmoked   int64
}

func (c counters) diff(before counters) counters {
	return counters{
		tobaccoMade:   c.tobaccoMade - before.tobaccoMade,
		paperMade:     c.paperMade - before.paperMade,
		matchMade:     c.matchMade - before.matchMade,
		tobaccoSmoked: c.tobaccoSmoked - before.tobaccoSmoked,
		paperSmoked:   c.paperSmoked - before.paperSmoked,
		matchSmoked:   c.matchSmoked - before.matchSmoked,
	}
}

type harness struct {
	agent *agent

	runDone     chan struct{}
	tobaccoDone chan struct{}
	paperDone   chan struct{}
	matchDone   chan struct{}

	tobaccoMade   atomic.Int64
	paperMade     atomic.Int64
	matchMade     atomic.Int64
	tobaccoSmoked atomic.Int64
	paperSmoked   atomic.Int64
	matchSmoked   atomic.Int64
}

func newHarness(ctx context.Context, newImpl func(agent Agent) CigaretteSmokers) *harness {
	a := newAgent()
	cs := newImpl(a)
	h := &harness{
		agent:       a,
		runDone:     make(chan struct{}),
		tobaccoDone: make(chan struct{}),
		paperDone:   make(chan struct{}),
		matchDone:   make(chan struct{}),
	}

	go func() {
		cs.Run(ctx)
		close(h.runDone)
	}()
	go func() {
		cs.SmokerWithTobacco(ctx, func() {
			h.tobaccoMade.Add(1)
		}, func() {
			h.tobaccoSmoked.Add(1)
		})
		close(h.tobaccoDone)
	}()
	go func() {
		cs.SmokerWithPaper(ctx, func() {
			h.paperMade.Add(1)
		}, func() {
			h.paperSmoked.Add(1)
		})
		close(h.paperDone)
	}()
	go func() {
		cs.SmokerWithMatch(ctx, func() {
			h.matchMade.Add(1)
		}, func() {
			h.matchSmoked.Add(1)
		})
		close(h.matchDone)
	}()

	return h
}

func (h *harness) snapshot() counters {
	return counters{
		tobaccoMade:   h.tobaccoMade.Load(),
		paperMade:     h.paperMade.Load(),
		matchMade:     h.matchMade.Load(),
		tobaccoSmoked: h.tobaccoSmoked.Load(),
		paperSmoked:   h.paperSmoked.Load(),
		matchSmoked:   h.matchSmoked.Load(),
	}
}

func runRound(t *testing.T, h *harness, pair ingredientPair) counters {
	t.Helper()

	before := h.snapshot()

	h.agent.beginRound(t, pair)
	synctest.Wait()
	h.agent.waitRoundComplete(t)
	synctest.Wait()

	return h.snapshot().diff(before)
}

func requireClosed(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
	default:
		t.Fatal(message)
	}
}

func verifySingleSmoker(t *testing.T, pair ingredientPair, got counters) {
	t.Helper()

	want := pair.expectedSmoker()

	if got.tobaccoMade != 0 && want != smokerWithTobacco {
		t.Fatalf("smoker with tobacco made %d cigarettes for %s", got.tobaccoMade, pair)
	}
	if got.paperMade != 0 && want != smokerWithPaper {
		t.Fatalf("smoker with paper made %d cigarettes for %s", got.paperMade, pair)
	}
	if got.matchMade != 0 && want != smokerWithMatch {
		t.Fatalf("smoker with match made %d cigarettes for %s", got.matchMade, pair)
	}
	if got.tobaccoSmoked != 0 && want != smokerWithTobacco {
		t.Fatalf("smoker with tobacco smoked %d cigarettes for %s", got.tobaccoSmoked, pair)
	}
	if got.paperSmoked != 0 && want != smokerWithPaper {
		t.Fatalf("smoker with paper smoked %d cigarettes for %s", got.paperSmoked, pair)
	}
	if got.matchSmoked != 0 && want != smokerWithMatch {
		t.Fatalf("smoker with match smoked %d cigarettes for %s", got.matchSmoked, pair)
	}

	switch want {
	case smokerWithTobacco:
		if got.tobaccoMade != 1 || got.tobaccoSmoked != 1 {
			t.Fatalf("smoker with tobacco got %+v, want one make and one smoke", got)
		}
	case smokerWithPaper:
		if got.paperMade != 1 || got.paperSmoked != 1 {
			t.Fatalf("smoker with paper got %+v, want one make and one smoke", got)
		}
	case smokerWithMatch:
		if got.matchMade != 1 || got.matchSmoked != 1 {
			t.Fatalf("smoker with match got %+v, want one make and one smoke", got)
		}
	}
}

func Run(t *testing.T, newImpl func(agent Agent) CigaretteSmokers) {
	t.Helper()

	t.Run("DispatchesCorrectSmoker", func(t *testing.T) {
		for _, pair := range []ingredientPair{tobaccoPaper, paperMatch, tobaccoMatch} {
			t.Run(pair.String(), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					h := newHarness(t.Context(), newImpl)
					synctest.Wait()

					verifySingleSmoker(t, pair, runRound(t, h, pair))
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h := newHarness(t.Context(), newImpl)
			synctest.Wait()

			sequence := []ingredientPair{
				tobaccoPaper,
				paperMatch,
				tobaccoMatch,
				tobaccoPaper,
				tobaccoMatch,
			}

			for _, pair := range sequence {
				verifySingleSmoker(t, pair, runRound(t, h, pair))
			}

			got := h.snapshot()
			if got.tobaccoSmoked != 1 {
				t.Fatalf("smoker with tobacco smoked %d cigarettes, want 1", got.tobaccoSmoked)
			}
			if got.paperSmoked != 2 {
				t.Fatalf("smoker with paper smoked %d cigarettes, want 2", got.paperSmoked)
			}
			if got.matchSmoked != 2 {
				t.Fatalf("smoker with match smoked %d cigarettes, want 2", got.matchSmoked)
			}
		})
	})

	t.Run("SignalsAgentBeforeSmoking", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			a := newAgent()
			cs := newImpl(a)

			ctx := t.Context()
			smokeStarted := make(chan struct{})
			releaseSmoke := make(chan struct{})
			tobaccoSmoked := make(chan struct{}, 1)

			go cs.Run(ctx)
			go cs.SmokerWithTobacco(ctx, func() {}, func() {
				tobaccoSmoked <- struct{}{}
			})
			go cs.SmokerWithPaper(ctx, func() {}, func() {})
			go cs.SmokerWithMatch(ctx, func() {}, func() {
				close(smokeStarted)
				<-releaseSmoke
			})

			synctest.Wait()

			select {
			case <-a.agent:
			default:
				t.Fatal("agent was not ready for the first round")
			}

			a.tobacco <- struct{}{}
			a.paper <- struct{}{}

			synctest.Wait()

			select {
			case <-smokeStarted:
			default:
				t.Fatal("smoker with match did not begin smoking")
			}

			select {
			case <-a.agent:
			default:
				t.Fatal("agent was not signaled before the first smoker finished smoking")
			}

			a.paper <- struct{}{}
			a.match <- struct{}{}

			synctest.Wait()

			select {
			case <-tobaccoSmoked:
			default:
				t.Fatal("second round did not complete while the first smoker was still smoking")
			}

			close(releaseSmoke)
			synctest.Wait()
		})
	})

	t.Run("CancelWhileIdle", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			h := newHarness(ctx, newImpl)
			synctest.Wait()

			cancel()
			synctest.Wait()

			requireClosed(t, h.runDone, "Run did not exit after cancellation")
			requireClosed(t, h.tobaccoDone, "SmokerWithTobacco did not exit after cancellation")
			requireClosed(t, h.paperDone, "SmokerWithPaper did not exit after cancellation")
			requireClosed(t, h.matchDone, "SmokerWithMatch did not exit after cancellation")
		})
	})

	t.Run("CancelWhileCollectingIngredients", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			h := newHarness(ctx, newImpl)
			synctest.Wait()

			select {
			case <-h.agent.agent:
			default:
				t.Fatal("agent was not ready for the first round")
			}
			h.agent.tobacco <- struct{}{}

			synctest.Wait()
			cancel()
			synctest.Wait()

			requireClosed(t, h.runDone, "Run did not exit after cancellation while waiting for the second ingredient")
			requireClosed(t, h.tobaccoDone, "SmokerWithTobacco did not exit after cancellation")
			requireClosed(t, h.paperDone, "SmokerWithPaper did not exit after cancellation")
			requireClosed(t, h.matchDone, "SmokerWithMatch did not exit after cancellation")
		})
	})
}

func Benchmark(b *testing.B, newImpl func(agent Agent) CigaretteSmokers) {
	b.Helper()
	b.ReportAllocs()

	for _, pair := range []ingredientPair{tobaccoPaper, paperMatch, tobaccoMatch} {
		b.Run(pair.String(), func(b *testing.B) {
			a := newAgent()
			cs := newImpl(a)
			ctx, cancel := context.WithCancel(b.Context())
			defer cancel()

			go cs.Run(ctx)
			go cs.SmokerWithTobacco(ctx, func() {}, func() {})
			go cs.SmokerWithPaper(ctx, func() {}, func() {})
			go cs.SmokerWithMatch(ctx, func() {}, func() {})

			for b.Loop() {
				<-a.agent

				switch pair {
				case tobaccoPaper:
					a.tobacco <- struct{}{}
					a.paper <- struct{}{}
				case paperMatch:
					a.paper <- struct{}{}
					a.match <- struct{}{}
				case tobaccoMatch:
					a.tobacco <- struct{}{}
					a.match <- struct{}{}
				}

				<-a.agent

				a.agent <- struct{}{}
			}
		})
	}
}
