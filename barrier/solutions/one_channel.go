// Package solutions contains implementations of the barrier problem.
package solutions

import (
	"context"
	"sync"
)

type OneChannel struct {
	parties int
	arrived int
	mu      sync.Mutex
	release chan struct{}
}

func NewOneChannel(parties int) *OneChannel {
	return &OneChannel{
		parties: parties,
		release: make(chan struct{}),
	}
}

func (b *OneChannel) Wait(ctx context.Context, phase1, phase2 func()) error {
	phase1()

	b.mu.Lock()
	release := b.release
	b.arrived++
	if b.arrived == b.parties {
		close(b.release)
		b.release = make(chan struct{})
		b.arrived = 0
	}
	b.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-release:
	}

	phase2()
	return nil
}
