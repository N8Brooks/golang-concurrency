package cigarettesmokers_test

import (
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"testing/synctest"

	cigarette_smokers_problem "github.com/N8Brooks/golang-concurrency/modified_cigarette_smokers_problem"
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
}

func newAgent() *agent {
	a := agent{
		tobacco: make(chan struct{}, 100),
		paper:   make(chan struct{}, 100),
		match:   make(chan struct{}, 100),
	}
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

func TestCigaretteSmokers(t *testing.T) {
	const numIterations = 100
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		a := newAgent()
		cs := cigarette_smokers_problem.NewCigaretteSmokers(t.Context(), a)

		var tobaccos, papers, matches atomic.Int64

		var smokerWithTobaccoCigarettes, smokerWithTobaccoActual atomic.Int64
		go cs.SmokerWithTobacco(func() {
			if papers.Add(-1) < 0 || matches.Add(-1) < 0 {
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

		var smokerWithPaperCigarettes, smokerWithPaperActual atomic.Int64
		go cs.SmokerWithPaper(func() {
			if tobaccos.Add(-1) < 0 || matches.Add(-1) < 0 {
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

		var smokerWithMatchCigarettes, smokerWithMatchActual atomic.Int64
		go cs.SmokerWithMatch(func() {
			if tobaccos.Add(-1) < 0 || papers.Add(-1) < 0 {
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

		go cs.Run()

		var smokerWithTobaccoExpected, smokerWithPaperExpected, smokerWithMatchExpected int64
		go func() {
			for range numIterations {
				expected := smokerID(rand.N(3) + 1)
				switch expected {
				case smokerWithTobacco:
					papers.Add(1)
					matches.Add(1)
					a.paper <- struct{}{}
					a.match <- struct{}{}
					smokerWithTobaccoExpected++
				case smokerWithPaper:
					tobaccos.Add(1)
					matches.Add(1)
					a.tobacco <- struct{}{}
					a.match <- struct{}{}
					smokerWithPaperExpected++
				case smokerWithMatch:
					tobaccos.Add(1)
					papers.Add(1)
					a.tobacco <- struct{}{}
					a.paper <- struct{}{}
					smokerWithMatchExpected++
				}
			}
		}()

		synctest.Wait()

		tobacco, p, m := tobaccos.Load(), papers.Load(), matches.Load()
		t.Logf("smoker with tobacco expected to smoke %d cigarettes, actually smoked %d cigarettes", smokerWithTobaccoExpected, smokerWithTobaccoActual.Load())
		t.Logf("smoker with paper expected to smoke %d cigarettes, actually smoked %d cigarettes", smokerWithPaperExpected, smokerWithPaperActual.Load())
		t.Logf("smoker with match expected to smoke %d cigarettes, actually smoked %d cigarettes", smokerWithMatchExpected, smokerWithMatchActual.Load())
		t.Logf("total tobacco supplied %d, remaining %d", smokerWithPaperExpected+smokerWithMatchExpected, tobacco)
		t.Logf("total paper supplied %d, remaining %d", smokerWithTobaccoExpected+smokerWithMatchExpected, p)
		t.Logf("total match supplied %d, remaining %d", smokerWithTobaccoExpected+smokerWithPaperExpected, m)
		if (tobacco > 0) && (p > 0) || (tobacco > 0) && (m > 0) || (p > 0) && (m > 0) {
			t.Error("expected all supplies to be used up, but some supplies were not used up")
		}
	})
}
