// Package solutions contains implementations of the barrier problem.
package solutions

import (
	"context"
	"sync"
)

type TwoChannels struct {
	parties int
	count   int
	mu      sync.Mutex
	first   chan struct{}
	second  chan struct{}
}

func NewTwoChannels(parties int) *TwoChannels {
	return &TwoChannels{
		parties: parties,
		first:   make(chan struct{}, parties),
		second:  make(chan struct{}, parties),
	}
}

func (b *TwoChannels) Wait(ctx context.Context, phase1, phase2 func()) error {
	phase1()

	b.mu.Lock()
	b.count++
	if b.count == b.parties {
		for range b.parties {
			b.first <- struct{}{}
		}
	}
	b.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.first:
	}

	b.mu.Lock()
	b.count--
	if b.count == 0 {
		for range b.parties {
			b.second <- struct{}{}
		}
	}
	b.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.second:
	}

	phase2()
	return nil
}
