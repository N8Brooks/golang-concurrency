package solutions

import (
	"context"
	"sync"
)

type syncMutexWaiter struct {
	ready chan struct{}
}

type SyncMutex struct {
	mu        sync.Mutex
	leaders   []*syncMutexWaiter
	followers []*syncMutexWaiter
}

func NewSyncMutex() *SyncMutex {
	return &SyncMutex{}
}

func (q *SyncMutex) Leader(ctx context.Context, leader, dance func()) error {
	leader()
	waiter := &syncMutexWaiter{ready: make(chan struct{})}

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
		q.leaders, removed = removeSyncMutexWaiter(q.leaders, waiter)
		q.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.ready
		dance()
		return nil
	}
}

func (q *SyncMutex) Follower(ctx context.Context, follower, dance func()) error {
	follower()
	waiter := &syncMutexWaiter{ready: make(chan struct{})}

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
		q.followers, removed = removeSyncMutexWaiter(q.followers, waiter)
		q.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.ready
		dance()
		return nil
	}
}

func removeSyncMutexWaiter(waiters []*syncMutexWaiter, target *syncMutexWaiter) ([]*syncMutexWaiter, bool) {
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
