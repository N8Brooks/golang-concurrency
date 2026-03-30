// Package solutions contains implementations of the cigarette smokers problem.
package solutions

import "context"

type Agent interface {
	Tobacco() chan struct{}
	Paper() chan struct{}
	Match() chan struct{}
	Signal()
}

type Channel struct {
	agent   Agent
	tobacco chan struct{}
	paper   chan struct{}
	match   chan struct{}
}

func NewChannel(agent Agent) *Channel {
	return &Channel{
		agent:   agent,
		tobacco: make(chan struct{}),
		paper:   make(chan struct{}),
		match:   make(chan struct{}),
	}
}

func (cs *Channel) Run(ctx context.Context) {
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

		switch {
		case isPaper && isMatch:
			cs.tobacco <- struct{}{}
		case isTobacco && isMatch:
			cs.paper <- struct{}{}
		case isTobacco && isPaper:
			cs.match <- struct{}{}
		default:
			panic("unexpected state: exactly two ingredients should be present")
		}
	}
}

func (cs *Channel) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
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

func (cs *Channel) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
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

func (cs *Channel) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
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
