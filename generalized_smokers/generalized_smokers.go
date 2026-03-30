//go:build challenge

// Package generalizedsmokers contains the challenge version of the generalized
// smokers problem.
//
// An agent places ingredients on the table without waiting for smokers, so
// multiple copies of tobacco, paper, and matches may accumulate. The solution
// must track available ingredients, wake the correct smoker whenever a complete
// pair is available, and keep running as more ingredients arrive.
package generalizedsmokers

import "context"

type Agent interface {
	Tobacco() chan struct{}
	Paper() chan struct{}
	Match() chan struct{}
}

type Smokers struct{}

func NewSmokers(a Agent) *Smokers {
	return &Smokers{}
}

func (s *Smokers) Run(ctx context.Context) {
	panic("unimplemented")
}

func (s *Smokers) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	panic("unimplemented")
}

func (s *Smokers) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	panic("unimplemented")
}

func (s *Smokers) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	panic("unimplemented")
}
