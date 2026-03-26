package cigarettesmokers_test

import (
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"testing/synctest"

	cigarette_smokers_problem "github.com/N8Brooks/golang-concurrency/cigarette_smokers_problem"
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

func (a *agent) SignalAgent() {
	a.agent <- struct{}{}
}

func TestCigaretteSmokers(t *testing.T) {
	const numIterations = 100
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		a := newAgent()
		cs := cigarette_smokers_problem.NewCigaretteSmokers(t.Context(), a)

		smokerWithTobaccoSupplies := make(chan struct{}, numIterations)
		smokerWithTobaccoCigarettes := make(chan struct{}, numIterations)
		smokerWithTobaccoSmoked := atomic.Int64{}
		go cs.SmokerWithTobacco(func() {
			select {
			case <-smokerWithTobaccoSupplies:
				t.Log("smoker with tobacco made a cigarette")
			default:
				t.Error("smoker with tobacco could not make a cigarette")
			}
			smokerWithTobaccoCigarettes <- struct{}{}
		}, func() {
			select {
			case <-smokerWithTobaccoCigarettes:
				t.Log("smoker with tobacco smoked a cigarette")
			default:
				t.Error("smoker with tobacco could not smoke a cigarette")
			}
			smokerWithTobaccoSmoked.Add(1)
		})

		smokerWithPaperSupplies := make(chan struct{}, numIterations)
		smokerWithPaperCigarettes := make(chan struct{}, numIterations)
		smokerWithPaperSmoked := atomic.Int64{}
		go cs.SmokerWithPaper(func() {
			select {
			case <-smokerWithPaperSupplies:
				t.Log("smoker with paper made a cigarette")
			default:
				t.Error("smoker with paper could not make a cigarette")
			}
			smokerWithPaperCigarettes <- struct{}{}
		}, func() {
			select {
			case <-smokerWithPaperCigarettes:
				t.Log("smoker with paper smoked a cigarette")
			default:
				t.Error("smoker with paper could not smoke a cigarette")
			}
			smokerWithPaperSmoked.Add(1)
		})

		smokerWithMatchSupplies := make(chan struct{}, numIterations)
		smokerWithMatchCigarettes := make(chan struct{}, numIterations)
		smokerWithMatchSmoked := atomic.Int64{}
		go cs.SmokerWithMatch(func() {
			select {
			case <-smokerWithMatchSupplies:
				t.Log("smoker with match made a cigarette")
			default:
				t.Error("smoker with match could not make a cigarette")
			}
			smokerWithMatchCigarettes <- struct{}{}
		}, func() {
			select {
			case <-smokerWithMatchCigarettes:
				t.Log("smoker with match smoked a cigarette")
			default:
				t.Error("smoker with match could not smoke a cigarette")
			}
			smokerWithMatchSmoked.Add(1)
		})

		go cs.Run()

		var smokerWithTobaccoExpected, smokerWithPaperExpected, smokerWithMatchExpected int64
		go func() {
			for range numIterations {
				<-a.agent
				expected := smokerID(rand.N(3) + 1)
				switch expected {
				case smokerWithTobacco:
					smokerWithTobaccoSupplies <- struct{}{}
					a.paper <- struct{}{}
					a.match <- struct{}{}
					smokerWithTobaccoExpected++
				case smokerWithPaper:
					smokerWithPaperSupplies <- struct{}{}
					a.tobacco <- struct{}{}
					a.match <- struct{}{}
					smokerWithPaperExpected++
				case smokerWithMatch:
					smokerWithMatchSupplies <- struct{}{}
					a.tobacco <- struct{}{}
					a.paper <- struct{}{}
					smokerWithMatchExpected++
				}
			}
		}()

		synctest.Wait()

		if smokerWithTobaccoSmoked.Load() != smokerWithTobaccoExpected {
			t.Errorf("expected smoker with tobacco to smoke %d cigarettes, but they smoked %d cigarettes", smokerWithTobaccoExpected, smokerWithTobaccoSmoked.Load())
		}
		if smokerWithPaperSmoked.Load() != smokerWithPaperExpected {
			t.Errorf("expected smoker with paper to smoke %d cigarettes, but they smoked %d cigarettes", smokerWithPaperExpected, smokerWithPaperSmoked.Load())
		}
		if smokerWithMatchSmoked.Load() != smokerWithMatchExpected {
			t.Errorf("expected smoker with match to smoke %d cigarettes, but they smoked %d cigarettes", smokerWithMatchExpected, smokerWithMatchSmoked.Load())
		}
	})
}
