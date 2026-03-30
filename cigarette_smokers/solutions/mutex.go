// Package solutions contains implementations of the cigarette smokers problem.
package solutions

import (
	"context"
	"sync"
)

type Mutex struct {
	agent Agent

	mu        sync.Mutex
	isTobacco bool
	isPaper   bool
	isMatch   bool

	tobaccoSem *countingSemaphore
	paperSem   *countingSemaphore
	matchSem   *countingSemaphore
}

func NewMutex(agent Agent) *Mutex {
	return &Mutex{
		agent:      agent,
		tobaccoSem: newCountingSemaphore(0),
		paperSem:   newCountingSemaphore(0),
		matchSem:   newCountingSemaphore(0),
	}
}

func (cs *Mutex) Run(ctx context.Context) {
	go cs.pusherTobacco(ctx)
	go cs.pusherPaper(ctx)
	go cs.pusherMatch(ctx)

	<-ctx.Done()
}

func (cs *Mutex) pusherTobacco(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.agent.Tobacco():
		}

		var target *countingSemaphore

		cs.mu.Lock()
		switch {
		case cs.isPaper:
			cs.isPaper = false
			target = cs.matchSem
		case cs.isMatch:
			cs.isMatch = false
			target = cs.paperSem
		default:
			cs.isTobacco = true
		}
		cs.mu.Unlock()

		if target != nil {
			target.Signal()
		}
	}
}

func (cs *Mutex) pusherPaper(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.agent.Paper():
		}

		var target *countingSemaphore

		cs.mu.Lock()
		switch {
		case cs.isTobacco:
			cs.isTobacco = false
			target = cs.matchSem
		case cs.isMatch:
			cs.isMatch = false
			target = cs.tobaccoSem
		default:
			cs.isPaper = true
		}
		cs.mu.Unlock()

		if target != nil {
			target.Signal()
		}
	}
}

func (cs *Mutex) pusherMatch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cs.agent.Match():
		}

		var target *countingSemaphore

		cs.mu.Lock()
		switch {
		case cs.isTobacco:
			cs.isTobacco = false
			target = cs.paperSem
		case cs.isPaper:
			cs.isPaper = false
			target = cs.tobaccoSem
		default:
			cs.isMatch = true
		}
		cs.mu.Unlock()

		if target != nil {
			target.Signal()
		}
	}
}

func (cs *Mutex) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if !cs.tobaccoSem.Wait(ctx) {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Mutex) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if !cs.paperSem.Wait(ctx) {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Mutex) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if !cs.matchSem.Wait(ctx) {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}
