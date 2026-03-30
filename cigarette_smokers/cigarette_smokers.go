//go:build challenge

// Package cigarettesmokers contains the challenge version of the cigarette
// smokers problem.
//
// An agent places two ingredients on the table, and exactly the smoker with
// the complementary ingredient should proceed. The agent interface is fixed, so
// the coordination logic belongs entirely to the smokers side.
package cigarettesmokers

import "context"

type Agent interface {
	Tobacco() chan struct{}
	Paper() chan struct{}
	Match() chan struct{}
	Signal()
}

type CigaretteSmokers struct {
	agent Agent
}

func NewCigaretteSmokers(agent Agent) *CigaretteSmokers {
	return &CigaretteSmokers{agent: agent}
}

func (cs *CigaretteSmokers) Run(ctx context.Context) {
	panic("unimplemented")
}

func (cs *CigaretteSmokers) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	panic("unimplemented")
}

func (cs *CigaretteSmokers) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	panic("unimplemented")
}

func (cs *CigaretteSmokers) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	panic("unimplemented")
}
