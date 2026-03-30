package cigarettesmokers

import (
	"context"
)

type Agent interface {
	// Tobacco returns a channel that signals when tobacco is available available.
	Tobacco() chan struct{}
	// Paper returns a channel that signals when paper is available.
	Paper() chan struct{}
	// Match returns a channel that signals when matches are available.
	Match() chan struct{}
	// Signal signals the agent to place two random items on the table.
	Signal()
}

type CigaretteSmokers struct {
	agent   Agent
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
}

func NewCigaretteSmokers(a Agent) *CigaretteSmokers {
	cs := CigaretteSmokers{
		agent:   a,
		tobacco: make(chan struct{}),
		paper:   make(chan struct{}),
		match:   make(chan struct{}),
	}
	return &cs
}

func (cs *CigaretteSmokers) Run(ctx context.Context) {
	for {
		var isTobacco, isPaper, isMatch bool

		select {
		case <-ctx.Done():
			return
		case <-cs.agent.Tobacco():
			isTobacco = true
		case <-cs.agent.Paper():
			isPaper = true
		case <-cs.agent.Match():
			isMatch = true
		}

		select {
		case <-ctx.Done():
			return
		case <-cs.agent.Tobacco():
			isTobacco = true
		case <-cs.agent.Paper():
			isPaper = true
		case <-cs.agent.Match():
			isMatch = true
		}

		if isPaper && isMatch {
			cs.tobacco <- struct{}{}
		} else if isTobacco && isMatch {
			cs.paper <- struct{}{}
		} else if isTobacco && isPaper {
			cs.match <- struct{}{}
		} else {
			panic("unexpected state: exactly two of tobacco, paper, and match should be true")
		}
	}
}

func (cs *CigaretteSmokers) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.tobacco:
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *CigaretteSmokers) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.paper:
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *CigaretteSmokers) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.match:
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}
