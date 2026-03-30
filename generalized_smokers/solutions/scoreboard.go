// Package solutions contains implementations of the generalized cigarette smokers problem.
package solutions

import "context"

type Agent interface {
	Tobacco() chan struct{}
	Paper() chan struct{}
	Match() chan struct{}
}

type Scoreboard struct {
	agent   Agent
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
}

func NewScoreboard(a Agent) *Scoreboard {
	return &Scoreboard{
		agent:   a,
		tobacco: make(chan struct{}),
		paper:   make(chan struct{}),
		match:   make(chan struct{}),
	}
}

func (cs *Scoreboard) Run(ctx context.Context) {
	var numTobacco, numPaper, numMatch int

	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.agent.Tobacco():
			numTobacco++
		case <-cs.agent.Paper():
			numPaper++
		case <-cs.agent.Match():
			numMatch++
		}

		if numPaper > 0 && numMatch > 0 {
			numPaper--
			numMatch--
			select {
			case cs.tobacco <- struct{}{}:
			case <-ctx.Done():
				return
			}
		} else if numTobacco > 0 && numMatch > 0 {
			numTobacco--
			numMatch--
			select {
			case cs.paper <- struct{}{}:
			case <-ctx.Done():
				return
			}
		} else if numTobacco > 0 && numPaper > 0 {
			numTobacco--
			numPaper--
			select {
			case cs.match <- struct{}{}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (cs *Scoreboard) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.tobacco:
		}
		makeCigarette()
		smoke()
	}
}

func (cs *Scoreboard) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.paper:
		}
		makeCigarette()
		smoke()
	}
}

func (cs *Scoreboard) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.match:
		}
		makeCigarette()
		smoke()
	}
}
