package cigarettesmokers

import (
	"context"
)

type Agent interface {
	// Tobacco returns a channel that signals when tobacco is available.
	Tobacco() chan struct{}
	// Paper returns a channel that signals when paper is available.
	Paper() chan struct{}
	// Match returns a channel that signals when matches are available.
	Match() chan struct{}
}

type CigaretteSmokers struct {
	ctx     context.Context
	agent   Agent
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
}

func NewCigaretteSmokers(ctx context.Context, a Agent) *CigaretteSmokers {
	cs := CigaretteSmokers{
		ctx:     ctx,
		agent:   a,
		tobacco: make(chan struct{}),
		paper:   make(chan struct{}),
		match:   make(chan struct{}),
	}
	return &cs
}

func (cs *CigaretteSmokers) Run() {
	for {
		var isTobacco, isPaper, isMatch bool

		select {
		case <-cs.ctx.Done():
			return
		case <-cs.agent.Tobacco():
			isTobacco = true
		case <-cs.agent.Paper():
			isPaper = true
		case <-cs.agent.Match():
			isMatch = true
		}

		select {
		case <-cs.ctx.Done():
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

func (cs *CigaretteSmokers) SmokerWithTobacco(smoke func()) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.tobacco:
		}
		smoke()
	}
}

func (cs *CigaretteSmokers) SmokerWithPaper(smoke func()) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.paper:
		}
		smoke()
	}
}

func (cs *CigaretteSmokers) SmokerWithMatch(smoke func()) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.match:
		}
		smoke()
	}
}
