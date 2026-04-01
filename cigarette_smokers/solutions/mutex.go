// Package solutions contains implementations of the cigarette smokers problem.
package solutions

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Mutex struct {
	agent Agent

	mu        sync.Mutex
	isTobacco bool
	isPaper   bool
	isMatch   bool

	tobaccoSem *semaphore.Weighted
	paperSem   *semaphore.Weighted
	matchSem   *semaphore.Weighted
}

func NewMutex(agent Agent) *Mutex {
	ctx := context.TODO()
	tobaccoSem := semaphore.NewWeighted(1)
	tobaccoSem.Acquire(ctx, 1)
	paperSem := semaphore.NewWeighted(1)
	paperSem.Acquire(ctx, 1)
	matchSem := semaphore.NewWeighted(1)
	matchSem.Acquire(ctx, 1)
	return &Mutex{
		agent:      agent,
		tobaccoSem: tobaccoSem,
		paperSem:   paperSem,
		matchSem:   matchSem,
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

		var target *semaphore.Weighted
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
			target.Release(1)
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

		var target *semaphore.Weighted

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
			target.Release(1)
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

		var target *semaphore.Weighted

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
			target.Release(1)
		}
	}
}

func (cs *Mutex) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if err := cs.tobaccoSem.Acquire(ctx, 1); err != nil {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Mutex) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if err := cs.paperSem.Acquire(ctx, 1); err != nil {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Mutex) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if err := cs.matchSem.Acquire(ctx, 1); err != nil {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}
