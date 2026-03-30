package solutions

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type semaphoreWaiter struct {
	ready   *semaphore.Weighted
	boarded *semaphore.Weighted
}

type Semaphore struct {
	capacity int
	mu       sync.Mutex
	waiting  []*semaphoreWaiter
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{capacity: capacity}
}

func (b *Semaphore) Rider(ctx context.Context, boardBus func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	waiter := newSemaphoreWaiter()

	b.mu.Lock()
	b.waiting = append(b.waiting, waiter)
	b.mu.Unlock()

	if err := waiter.ready.Acquire(ctx, 1); err != nil {
		b.mu.Lock()
		removed := removeSemaphoreWaiter(&b.waiting, waiter)
		b.mu.Unlock()
		if removed {
			return err
		}
		if err := waiter.ready.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
	}

	boardBus()
	waiter.boarded.Release(1)
	return nil
}

func (b *Semaphore) Bus(ctx context.Context, depart func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.mu.Lock()
	n := min(len(b.waiting), b.capacity)
	if n == 0 {
		b.mu.Unlock()
		depart()
		return nil
	}
	selected := append([]*semaphoreWaiter(nil), b.waiting[:n]...)
	b.waiting = b.waiting[n:]
	b.mu.Unlock()

	for _, waiter := range selected {
		waiter.ready.Release(1)
	}
	for _, waiter := range selected {
		if err := waiter.boarded.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
	}

	depart()
	return nil
}

func newSemaphoreWaiter() *semaphoreWaiter {
	ready := semaphore.NewWeighted(1)
	boarded := semaphore.NewWeighted(1)
	if err := ready.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	if err := boarded.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	return &semaphoreWaiter{
		ready:   ready,
		boarded: boarded,
	}
}

func removeSemaphoreWaiter(waiters *[]*semaphoreWaiter, target *semaphoreWaiter) bool {
	for i, waiter := range *waiters {
		if waiter != target {
			continue
		}
		copy((*waiters)[i:], (*waiters)[i+1:])
		last := len(*waiters) - 1
		(*waiters)[last] = nil
		*waiters = (*waiters)[:last]
		return true
	}
	return false
}
