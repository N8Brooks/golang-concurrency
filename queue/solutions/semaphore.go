package solutions

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type semaphoreWaiter struct {
	ready *semaphore.Weighted
}

type Semaphore struct {
	mu        sync.Mutex
	leaders   []*semaphoreWaiter
	followers []*semaphoreWaiter
}

func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

func (q *Semaphore) Leader(ctx context.Context, leader, dance func()) error {
	leader()
	waiter := newSemaphoreWaiter()

	q.mu.Lock()
	if len(q.followers) > 0 {
		follower := q.followers[0]
		q.followers = q.followers[1:]
		q.mu.Unlock()
		follower.ready.Release(1)
		dance()
		return nil
	}
	q.leaders = append(q.leaders, waiter)
	q.mu.Unlock()

	if err := waiter.ready.Acquire(ctx, 1); err != nil {
		q.mu.Lock()
		var removed bool
		q.leaders, removed = removeSemaphoreWaiter(q.leaders, waiter)
		q.mu.Unlock()
		if removed {
			return err
		}
		if err := waiter.ready.Acquire(context.Background(), 1); err != nil {
			panic("semaphore waiter matched but did not wake")
		}
	}

	dance()
	return nil
}

func (q *Semaphore) Follower(ctx context.Context, follower, dance func()) error {
	follower()
	waiter := newSemaphoreWaiter()

	q.mu.Lock()
	if len(q.leaders) > 0 {
		leader := q.leaders[0]
		q.leaders = q.leaders[1:]
		q.mu.Unlock()
		leader.ready.Release(1)
		dance()
		return nil
	}
	q.followers = append(q.followers, waiter)
	q.mu.Unlock()

	if err := waiter.ready.Acquire(ctx, 1); err != nil {
		q.mu.Lock()
		var removed bool
		q.followers, removed = removeSemaphoreWaiter(q.followers, waiter)
		q.mu.Unlock()
		if removed {
			return err
		}
		if err := waiter.ready.Acquire(context.Background(), 1); err != nil {
			panic("semaphore waiter matched but did not wake")
		}
	}

	dance()
	return nil
}

func newSemaphoreWaiter() *semaphoreWaiter {
	ready := semaphore.NewWeighted(1)
	if err := ready.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	return &semaphoreWaiter{ready: ready}
}

func removeSemaphoreWaiter(waiters []*semaphoreWaiter, target *semaphoreWaiter) ([]*semaphoreWaiter, bool) {
	for i, waiter := range waiters {
		if waiter != target {
			continue
		}
		copy(waiters[i:], waiters[i+1:])
		waiters[len(waiters)-1] = nil
		return waiters[:len(waiters)-1], true
	}
	return waiters, false
}
