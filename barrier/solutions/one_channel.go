// Package solutions contains implementations of the barrier problem.
package solutions

import (
	"context"
	"sync"
)

type OneChannel struct {
	parties   int
	count     int
	mu        sync.Mutex
	turnstile chan struct{}
}

func NewOneChannel(parties int) *OneChannel {
	return &OneChannel{
		parties:   parties,
		turnstile: make(chan struct{}, 1),
	}
}

func (b *OneChannel) Wait(ctx context.Context, phase1, phase2 func()) error {
	phase1()

	b.mu.Lock()
	b.count++
	last := b.count == b.parties
	b.mu.Unlock()

	if last {
		b.turnstile <- struct{}{}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.turnstile:
	}

	b.turnstile <- struct{}{}

	phase2()
	return nil
}
