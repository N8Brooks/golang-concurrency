// Package solutions contains implementations of the barrier problem.
package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	parties int
	count   int
	mutex   *semaphore.Weighted
	barrier *semaphore.Weighted
}

func NewSemaphore(parties int) *Semaphore {
	barrier := semaphore.NewWeighted(1)
	barrier.Acquire(context.Background(), 1)

	return &Semaphore{
		parties: parties,
		mutex:   semaphore.NewWeighted(1),
		barrier: barrier,
	}
}

func (b *Semaphore) Wait(ctx context.Context, phase1, phase2 func()) error {
	phase1()

	if err := b.mutex.Acquire(ctx, 1); err != nil {
		return err
	}

	b.count++
	last := b.count == b.parties
	b.mutex.Release(1)

	if last {
		b.barrier.Release(1)
	}

	if err := b.barrier.Acquire(ctx, 1); err != nil {
		return err
	}

	b.barrier.Release(1)

	phase2()
	return nil
}
