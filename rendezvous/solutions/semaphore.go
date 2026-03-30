// Package solutions contains implementations of the rendezvous problem.
package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	a *semaphore.Weighted
	b *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	a := semaphore.NewWeighted(1)
	b := semaphore.NewWeighted(1)
	a.Acquire(context.Background(), 1)
	b.Acquire(context.Background(), 1)
	return &Semaphore{
		a: a,
		b: b,
	}
}

func (r *Semaphore) A(ctx context.Context, a1, a2 func()) {
	a1()
	r.a.Release(1)
	if err := r.b.Acquire(ctx, 1); err != nil {
		return
	}
	a2()
}

func (r *Semaphore) B(ctx context.Context, b1, b2 func()) {
	b1()
	r.b.Release(1)
	if err := r.a.Acquire(ctx, 1); err != nil {
		return
	}
	b2()
}
