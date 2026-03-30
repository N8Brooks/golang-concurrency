// Package solutions contains implementations of the queue problem.
package solutions

import (
	"context"
	"sync"
)

type channelWaiter struct {
	ready chan struct{}
}

type Channel struct {
	mu        sync.Mutex
	leaders   []*channelWaiter
	followers []*channelWaiter
}

func NewChannel() *Channel {
	return &Channel{}
}

func (q *Channel) Leader(ctx context.Context, leader, dance func()) error {
	leader()
	waiter := &channelWaiter{ready: make(chan struct{})}

	q.mu.Lock()
	if len(q.followers) > 0 {
		follower := q.followers[0]
		q.followers = q.followers[1:]
		q.mu.Unlock()
		close(follower.ready)
		dance()
		return nil
	}
	q.leaders = append(q.leaders, waiter)
	q.mu.Unlock()

	select {
	case <-waiter.ready:
		dance()
		return nil
	case <-ctx.Done():
		q.mu.Lock()
		var removed bool
		q.leaders, removed = removeChannelWaiter(q.leaders, waiter)
		q.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.ready
		dance()
		return nil
	}
}

func (q *Channel) Follower(ctx context.Context, follower, dance func()) error {
	follower()
	waiter := &channelWaiter{ready: make(chan struct{})}

	q.mu.Lock()
	if len(q.leaders) > 0 {
		leader := q.leaders[0]
		q.leaders = q.leaders[1:]
		q.mu.Unlock()
		close(leader.ready)
		dance()
		return nil
	}
	q.followers = append(q.followers, waiter)
	q.mu.Unlock()

	select {
	case <-waiter.ready:
		dance()
		return nil
	case <-ctx.Done():
		q.mu.Lock()
		var removed bool
		q.followers, removed = removeChannelWaiter(q.followers, waiter)
		q.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.ready
		dance()
		return nil
	}
}

func removeChannelWaiter(waiters []*channelWaiter, target *channelWaiter) ([]*channelWaiter, bool) {
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
