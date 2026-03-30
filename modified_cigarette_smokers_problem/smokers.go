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
	var numTobacco, numPaper, numMatch int
	for {
		select {
		case <-cs.ctx.Done():
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
			case <-cs.ctx.Done():
				return
			}
		} else if numTobacco > 0 && numMatch > 0 {
			numTobacco--
			numMatch--
			select {
			case cs.paper <- struct{}{}:
			case <-cs.ctx.Done():
				return
			}
		} else if numTobacco > 0 && numPaper > 0 {
			numTobacco--
			numPaper--
			select {
			case cs.match <- struct{}{}:
			case <-cs.ctx.Done():
				return
			}
		}
	}
}

func (cs *CigaretteSmokers) SmokerWithTobacco(makeCigarette, smoke func()) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.tobacco:
		}
		makeCigarette()
		smoke()
	}
}

func (cs *CigaretteSmokers) SmokerWithPaper(makeCigarette, smoke func()) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.paper:
		}
		makeCigarette()
		smoke()
	}
}

func (cs *CigaretteSmokers) SmokerWithMatch(makeCigarette, smoke func()) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.match:
		}
		makeCigarette()
		smoke()
	}
}
