// Package solutions contains implementations of the barrier problem.
package solutions

import (
	"context"
)

type TwoChannels struct {
	parties   int
	arrivals  chan struct{}
	turnstile chan struct{}
}

func NewTwoChannels(parties int) *TwoChannels {
	return &TwoChannels{
		parties:   parties,
		arrivals:  make(chan struct{}, parties),
		turnstile: make(chan struct{}, 1),
	}
}

func (b *TwoChannels) Wait(ctx context.Context, phase1, phase2 func()) error {
	phase1()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.arrivals <- struct{}{}:
	}

	if len(b.arrivals) == b.parties {
		select {
		case b.turnstile <- struct{}{}:
		default:
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.turnstile:
	}

	select {
	case b.turnstile <- struct{}{}:
	default:
	}

	phase2()
	return nil
}
