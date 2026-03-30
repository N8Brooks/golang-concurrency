package solutions

import (
	"context"
	"sync"
)

type syncMutexPair struct {
	followerDone chan struct{}
}

type syncMutexWaiter struct {
	pair chan *syncMutexPair
}

type SyncMutex struct {
	mu        sync.Mutex
	active    bool
	leaders   []*syncMutexWaiter
	followers []*syncMutexWaiter
}

func NewSyncMutex() *SyncMutex {
	return &SyncMutex{}
}

func (q *SyncMutex) Leader(ctx context.Context, leader, dance func()) error {
	leader()
	waiter := &syncMutexWaiter{pair: make(chan *syncMutexPair, 1)}

	q.mu.Lock()
	if !q.active && len(q.followers) > 0 {
		follower := q.followers[0]
		q.followers = q.followers[1:]
		pair := &syncMutexPair{followerDone: make(chan struct{})}
		q.active = true
		q.mu.Unlock()
		follower.pair <- pair
		dance()
		<-pair.followerDone
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
	<-pair.followerDone
	q.finishLeader()
	return nil
}

func (q *SyncMutex) Follower(ctx context.Context, follower, dance func()) error {
	follower()
	waiter := &syncMutexWaiter{pair: make(chan *syncMutexPair, 1)}

	q.mu.Lock()
	if !q.active && len(q.leaders) > 0 {
		leader := q.leaders[0]
		q.leaders = q.leaders[1:]
		pair := &syncMutexPair{followerDone: make(chan struct{})}
		q.active = true
		q.mu.Unlock()
		leader.pair <- pair
		dance()
		close(pair.followerDone)
		return nil
	}
	q.followers = append(q.followers, waiter)
	q.mu.Unlock()

	pair, err := q.waitFollower(ctx, waiter)
	if err != nil {
		return err
	}
	dance()
	close(pair.followerDone)
	return nil
}

func (q *SyncMutex) waitLeader(ctx context.Context, waiter *syncMutexWaiter) (*syncMutexPair, error) {
	select {
	case pair := <-waiter.pair:
		return pair, nil
	case <-ctx.Done():
		q.mu.Lock()
		var removed bool
		q.leaders, removed = removeSyncMutexWaiter(q.leaders, waiter)
		q.mu.Unlock()
		if removed {
			return nil, ctx.Err()
		}
		return <-waiter.pair, nil
	}
}

func (q *SyncMutex) waitFollower(ctx context.Context, waiter *syncMutexWaiter) (*syncMutexPair, error) {
	select {
	case pair := <-waiter.pair:
		return pair, nil
	case <-ctx.Done():
		q.mu.Lock()
		var removed bool
		q.followers, removed = removeSyncMutexWaiter(q.followers, waiter)
		q.mu.Unlock()
		if removed {
			return nil, ctx.Err()
		}
		return <-waiter.pair, nil
	}
}

func (q *SyncMutex) finishLeader() {
	q.mu.Lock()
	q.active = false
	q.activateNextLocked()
	q.mu.Unlock()
}

func (q *SyncMutex) activateNextLocked() {
	if q.active || len(q.leaders) == 0 || len(q.followers) == 0 {
		return
	}

	leader := q.leaders[0]
	follower := q.followers[0]
	q.leaders = q.leaders[1:]
	q.followers = q.followers[1:]

	pair := &syncMutexPair{followerDone: make(chan struct{})}
	q.active = true
	leader.pair <- pair
	follower.pair <- pair
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
