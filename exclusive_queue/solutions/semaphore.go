package solutions

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type semaphorePair struct {
	followerDone *semaphore.Weighted
}

type semaphoreWaiter struct {
	ready *semaphore.Weighted
	pair  *semaphorePair
}

type Semaphore struct {
	mu        sync.Mutex
	active    bool
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
	if !q.active && len(q.followers) > 0 {
		follower := q.followers[0]
		q.followers = q.followers[1:]
		pair := newSemaphorePair()
		follower.pair = pair
		q.active = true
		q.mu.Unlock()
		follower.ready.Release(1)
		dance()
		if err := pair.followerDone.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
		q.finishLeader()
		return nil
	}
	q.leaders = append(q.leaders, waiter)
	q.mu.Unlock()

	pair, err := q.waitLeader(ctx, waiter)
	if err != nil {
		return err
	}
	dance()
	if err := pair.followerDone.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	q.finishLeader()
	return nil
}

func (q *Semaphore) Follower(ctx context.Context, follower, dance func()) error {
	follower()
	waiter := newSemaphoreWaiter()

	q.mu.Lock()
	if !q.active && len(q.leaders) > 0 {
		leader := q.leaders[0]
		q.leaders = q.leaders[1:]
		pair := newSemaphorePair()
		leader.pair = pair
		q.active = true
		q.mu.Unlock()
		leader.ready.Release(1)
		dance()
		pair.followerDone.Release(1)
		return nil
	}
	q.followers = append(q.followers, waiter)
	q.mu.Unlock()

	pair, err := q.waitFollower(ctx, waiter)
	if err != nil {
		return err
	}
	dance()
	pair.followerDone.Release(1)
	return nil
}

func (q *Semaphore) waitLeader(ctx context.Context, waiter *semaphoreWaiter) (*semaphorePair, error) {
	if err := waiter.ready.Acquire(ctx, 1); err != nil {
		q.mu.Lock()
		var removed bool
		q.leaders, removed = removeSemaphoreWaiter(q.leaders, waiter)
		q.mu.Unlock()
		if removed {
			return nil, err
		}
		if err := waiter.ready.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
	}
	return waiter.pair, nil
}

func (q *Semaphore) waitFollower(ctx context.Context, waiter *semaphoreWaiter) (*semaphorePair, error) {
	if err := waiter.ready.Acquire(ctx, 1); err != nil {
		q.mu.Lock()
		var removed bool
		q.followers, removed = removeSemaphoreWaiter(q.followers, waiter)
		q.mu.Unlock()
		if removed {
			return nil, err
		}
		if err := waiter.ready.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
	}
	return waiter.pair, nil
}

func (q *Semaphore) finishLeader() {
	q.mu.Lock()
	q.active = false
	q.activateNextLocked()
	q.mu.Unlock()
}

func (q *Semaphore) activateNextLocked() {
	if q.active || len(q.leaders) == 0 || len(q.followers) == 0 {
		return
	}

	leader := q.leaders[0]
	follower := q.followers[0]
	q.leaders = q.leaders[1:]
	q.followers = q.followers[1:]

	pair := newSemaphorePair()
	leader.pair = pair
	follower.pair = pair
	q.active = true
	leader.ready.Release(1)
	follower.ready.Release(1)
}

func newSemaphoreWaiter() *semaphoreWaiter {
	ready := semaphore.NewWeighted(1)
	if err := ready.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	return &semaphoreWaiter{ready: ready}
}

func newSemaphorePair() *semaphorePair {
	done := semaphore.NewWeighted(1)
	if err := done.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	return &semaphorePair{followerDone: done}
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
