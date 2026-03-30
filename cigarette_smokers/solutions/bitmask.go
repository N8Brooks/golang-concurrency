// Package solutions contains implementations of the cigarette smokers problem.
package solutions

import (
	"context"
	"sync"
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
	smokers   [8]*countingSemaphore
}

func NewBitmask(agent Agent) *Bitmask {
	cs := &Bitmask{agent: agent}
	cs.smokers[tobaccoBit|paperBit] = newCountingSemaphore(0)
	cs.smokers[tobaccoBit|matchBit] = newCountingSemaphore(0)
	cs.smokers[paperBit|matchBit] = newCountingSemaphore(0)
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

		var target *countingSemaphore

		cs.mu.Lock()
		cs.available |= bit
		switch cs.available {
		case tobaccoBit | paperBit, tobaccoBit | matchBit, paperBit | matchBit:
			target = cs.smokers[cs.available]
			cs.available = 0
		}
		cs.mu.Unlock()

		if target != nil {
			target.Signal()
		}
	}
}

func (cs *Bitmask) SmokerWithTobacco(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if !cs.smokers[paperBit|matchBit].Wait(ctx) {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Bitmask) SmokerWithPaper(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if !cs.smokers[tobaccoBit|matchBit].Wait(ctx) {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}

func (cs *Bitmask) SmokerWithMatch(ctx context.Context, makeCigarette, smoke func()) {
	for {
		if !cs.smokers[tobaccoBit|paperBit].Wait(ctx) {
			return
		}
		makeCigarette()
		cs.agent.Signal()
		smoke()
	}
}
