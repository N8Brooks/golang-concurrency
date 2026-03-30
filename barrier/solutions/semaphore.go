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
	first   *semaphore.Weighted
	second  *semaphore.Weighted
}

func NewSemaphore(parties int) *Semaphore {
	first := semaphore.NewWeighted(int64(parties))
	second := semaphore.NewWeighted(int64(parties))
	first.Acquire(context.Background(), int64(parties))
	second.Acquire(context.Background(), int64(parties))

	return &Semaphore{
		parties: parties,
		first:   first,
		second:  second,
		mutex:   semaphore.NewWeighted(1),
	}
}

func (b *Semaphore) Wait(ctx context.Context, phase1, phase2 func()) error {
	phase1()

	if err := b.mutex.Acquire(ctx, 1); err != nil {
		return err
	}

	b.count++
	if b.count == b.parties {
		b.first.Release(int64(b.parties))
	}
	b.mutex.Release(1)

	if err := b.first.Acquire(ctx, 1); err != nil {
		return err
	}

	if err := b.mutex.Acquire(ctx, 1); err != nil {
		return err
	}

	b.count--
	if b.count == 0 {
		b.second.Release(int64(b.parties))
	}
	b.mutex.Release(1)

	if err := b.second.Acquire(ctx, 1); err != nil {
		return err
	}

	phase2()
	return nil
}
