package cigarettesmokers_test

import (
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"testing/synctest"

	cigarette_smokers_problem "github.com/N8Brooks/golang-concurrency/cigarette_smokers"
)

type smokerID int

const (
	smokerWithTobacco smokerID = iota + 1
	smokerWithPaper
	smokerWithMatch
)

type agent struct {
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
	agent   chan struct{}
}

func newAgent() *agent {
	a := agent{
		tobacco: make(chan struct{}),
		paper:   make(chan struct{}),
		match:   make(chan struct{}),
		agent:   make(chan struct{}, 1),
	}
	a.agent <- struct{}{}
	return &a
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

func TestCigaretteSmokers(t *testing.T) {
	const numIterations = 100
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		a := newAgent()
		ctx := t.Context()
		cs := cigarette_smokers_problem.NewCigaretteSmokers(a)

		var smokerWithTobaccoSupplies, smokerWithTobaccoCigarettes, smokerWithTobaccoActual atomic.Int64
		go cs.SmokerWithTobacco(ctx, func() {
			if smokerWithTobaccoSupplies.Add(-1) < 0 {
				t.Error("smoker with tobacco could not make a cigarette")
			} else {
				t.Log("smoker with tobacco made a cigarette")
			}
			smokerWithTobaccoCigarettes.Add(1)
		}, func() {
			if smokerWithTobaccoCigarettes.Add(-1) < 0 {
				t.Error("smoker with tobacco could not smoke a cigarette")
			} else {
				t.Log("smoker with tobacco smoked a cigarette")
			}
			smokerWithTobaccoActual.Add(1)
		})

		var smokerWithPaperSupplies, smokerWithPaperCigarettes, smokerWithPaperActual atomic.Int64
		go cs.SmokerWithPaper(ctx, func() {
			if smokerWithPaperSupplies.Add(-1) < 0 {
				t.Error("smoker with paper could not make a cigarette")
			} else {
				t.Log("smoker with paper made a cigarette")
			}
			smokerWithPaperCigarettes.Add(1)
		}, func() {
			if smokerWithPaperCigarettes.Add(-1) < 0 {
				t.Error("smoker with paper could not smoke a cigarette")
			} else {
				t.Log("smoker with paper smoked a cigarette")
			}
			smokerWithPaperActual.Add(1)
		})

		var smokerWithMatchSupplies, smokerWithMatchCigarettes, smokerWithMatchActual atomic.Int64
		go cs.SmokerWithMatch(ctx, func() {
			if smokerWithMatchSupplies.Add(-1) < 0 {
				t.Error("smoker with match could not make a cigarette")
			} else {
				t.Log("smoker with match made a cigarette")
			}
			smokerWithMatchCigarettes.Add(1)
		}, func() {
			if smokerWithMatchCigarettes.Add(-1) < 0 {
				t.Error("smoker with match could not smoke a cigarette")
			} else {
				t.Log("smoker with match smoked a cigarette")
			}
			smokerWithMatchActual.Add(1)
		})

		go cs.Run(ctx)

		var smokerWithTobaccoExpected, smokerWithPaperExpected, smokerWithMatchExpected int64
		go func() {
			for range numIterations {
				<-a.agent
				expected := smokerID(rand.N(3) + 1)
				switch expected {
				case smokerWithTobacco:
					smokerWithTobaccoSupplies.Add(1)
					a.paper <- struct{}{}
					a.match <- struct{}{}
					smokerWithTobaccoExpected++
				case smokerWithPaper:
					smokerWithPaperSupplies.Add(1)
					a.tobacco <- struct{}{}
					a.match <- struct{}{}
					smokerWithPaperExpected++
				case smokerWithMatch:
					smokerWithMatchSupplies.Add(1)
					a.tobacco <- struct{}{}
					a.paper <- struct{}{}
					smokerWithMatchExpected++
				}
			}
		}()

		synctest.Wait()

		if smokerWithTobaccoActual.Load() != smokerWithTobaccoExpected {
			t.Errorf("expected smoker with tobacco to smoke %d cigarettes, but they smoked %d cigarettes", smokerWithTobaccoExpected, smokerWithTobaccoActual.Load())
		}
		if smokerWithPaperActual.Load() != smokerWithPaperExpected {
			t.Errorf("expected smoker with paper to smoke %d cigarettes, but they smoked %d cigarettes", smokerWithPaperExpected, smokerWithPaperActual.Load())
		}
		if smokerWithMatchActual.Load() != smokerWithMatchExpected {
			t.Errorf("expected smoker with match to smoke %d cigarettes, but they smoked %d cigarettes", smokerWithMatchExpected, smokerWithMatchActual.Load())
		}
	})
}
