// Package solutions contains implementations of the cigarette smokers problem.
package solutions

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

const (
	tobaccoBit uint8 = 1 << iota
	paperBit
	matchBit
)

type Bitmask struct {
	agent Agent

	mu        sync.Mutex
	available uint8
	smokers   [8]*semaphore.Weighted
}

func NewBitmask(agent Agent) *Bitmask {
	cs := &Bitmask{agent: agent}
	ctx := context.TODO()
	tobaccoSem := semaphore.NewWeighted(1)
	tobaccoSem.Acquire(ctx, 1)
	cs.smokers[paperBit|matchBit] = tobaccoSem
	paperSem := semaphore.NewWeighted(1)
	paperSem.Acquire(ctx, 1)
	cs.smokers[tobaccoBit|matchBit] = paperSem
	matchSem := semaphore.NewWeighted(1)
	matchSem.Acquire(ctx, 1)
	cs.smokers[tobaccoBit|paperBit] = matchSem
	return cs
}

func (cs *Bitmask) Run(ctx context.Context) {
	go cs.pusher(ctx, cs.agent.Tobacco(), tobaccoBit)
	go cs.pusher(ctx, cs.agent.Paper(), paperBit)
	go cs.pusher(ctx, cs.agent.Match(), matchBit)

	<-ctx.Done()
}

func (cs *Bitmask) pusher(ctx context.Context, ingredient chan struct{}, bit uint8) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ingredient:
		}

		var target *semaphore.Weighted
		cs.mu.Lock()
		cs.available |= bit
		switch cs.available {
		case tobaccoBit | paperBit, tobaccoBit | matchBit, paperBit | matchBit:
			target = cs.smokers[cs.available]
			cs.available = 0
		}
		cs.mu.Unlock()

		if target != nil {
			target.Release(1)
		}
	}
}

func (cs *Bitmask) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if err := cs.smokers[paperBit|matchBit].Acquire(ctx, 1); err != nil {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Bitmask) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if err := cs.smokers[tobaccoBit|matchBit].Acquire(ctx, 1); err != nil {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Bitmask) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if err := cs.smokers[tobaccoBit|paperBit].Acquire(ctx, 1); err != nil {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}
