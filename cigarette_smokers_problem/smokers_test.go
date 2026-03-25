package cigarettesmokers_test

import (
	"math/rand/v2"
	"testing"
	"testing/synctest"

	cigarette_smokers_problem "github.com/N8Brooks/golang-concurrency/cigarette_smokers_problem"
)

type smokerID int

const (
	SmokerWithTobacco smokerID = iota + 1
	SmokerWithPaper
	SmokerWithMatch
)

type Agent struct {
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
}

func (a *Agent) Tobacco() chan struct{} {
	return a.tobacco
}

func (a *Agent) Paper() chan struct{} {
	return a.paper
}

func (a *Agent) Match() chan struct{} {
	return a.match
}

func TestCigaretteSmokers(t *testing.T) {
	const numIterations = 100
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		a := &Agent{
			match:   make(chan struct{}),
			tobacco: make(chan struct{}),
			paper:   make(chan struct{}),
		}
		cs := cigarette_smokers_problem.NewCigaretteSmokers(t.Context(), a)
		smokerIDs := make(chan smokerID)

		go cs.SmokerWithTobacco(func() {
			smokerIDs <- SmokerWithTobacco
		})

		go cs.SmokerWithPaper(func() {
			smokerIDs <- SmokerWithPaper
		})

		go cs.SmokerWithMatch(func() {
			smokerIDs <- SmokerWithMatch
		})

		go cs.Start()

		for range numIterations {
			expected := smokerID(rand.N(3) + 1)
			switch expected {
			case SmokerWithTobacco:
				a.paper <- struct{}{}
				a.match <- struct{}{}
			case SmokerWithPaper:
				a.tobacco <- struct{}{}
				a.match <- struct{}{}
			case SmokerWithMatch:
				a.tobacco <- struct{}{}
				a.paper <- struct{}{}
			}
			actual := <-smokerIDs
			if actual != expected {
				t.Fatalf("smokerId = %d, want %d", actual, expected)
			}
		}
	})
}
