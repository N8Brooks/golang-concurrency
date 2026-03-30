// Package solutions contains implementations of the exclusive queue problem.
package solutions

import (
	"context"
	"sync"
)

type channelPair struct {
	followerDone chan struct{}
}

type channelWaiter struct {
	pair chan *channelPair
}

type Channel struct {
	mu        sync.Mutex
	active    bool
	leaders   []*channelWaiter
	followers []*channelWaiter
}

func NewChannel() *Channel {
	return &Channel{}
}

func (q *Channel) Leader(ctx context.Context, leader, dance func()) error {
	leader()
	waiter := &channelWaiter{pair: make(chan *channelPair, 1)}

	q.mu.Lock()
	if !q.active && len(q.followers) > 0 {
		follower := q.followers[0]
		q.followers = q.followers[1:]
		pair := &channelPair{followerDone: make(chan struct{})}
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

func (q *Channel) Follower(ctx context.Context, follower, dance func()) error {
	follower()
	waiter := &channelWaiter{pair: make(chan *channelPair, 1)}

	q.mu.Lock()
	if !q.active && len(q.leaders) > 0 {
		leader := q.leaders[0]
		q.leaders = q.leaders[1:]
		pair := &channelPair{followerDone: make(chan struct{})}
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

func (q *Channel) waitLeader(ctx context.Context, waiter *channelWaiter) (*channelPair, error) {
	select {
	case pair := <-waiter.pair:
		return pair, nil
	case <-ctx.Done():
		q.mu.Lock()
		var removed bool
		q.leaders, removed = removeChannelWaiter(q.leaders, waiter)
		q.mu.Unlock()
		if removed {
			return nil, ctx.Err()
		}
		return <-waiter.pair, nil
	}
}

func (q *Channel) waitFollower(ctx context.Context, waiter *channelWaiter) (*channelPair, error) {
	select {
	case pair := <-waiter.pair:
		return pair, nil
	case <-ctx.Done():
		q.mu.Lock()
		var removed bool
		q.followers, removed = removeChannelWaiter(q.followers, waiter)
		q.mu.Unlock()
		if removed {
			return nil, ctx.Err()
		}
		return <-waiter.pair, nil
	}
}

func (q *Channel) finishLeader() {
	q.mu.Lock()
	q.active = false
	q.activateNextLocked()
	q.mu.Unlock()
}

func (q *Channel) activateNextLocked() {
	if q.active || len(q.leaders) == 0 || len(q.followers) == 0 {
		return
	}

	leader := q.leaders[0]
	follower := q.followers[0]
	q.leaders = q.leaders[1:]
	q.followers = q.followers[1:]

	pair := &channelPair{followerDone: make(chan struct{})}
	q.active = true
	leader.pair <- pair
	follower.pair <- pair
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
