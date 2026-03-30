// Package solutions contains implementations of the no-starve mutex problem.
package solutions

import (
	"context"
	"sync"
)

type channelWaiter struct {
	ready chan struct{}
}

type Channel struct {
	mu      sync.Mutex
	locked  bool
	waiters []*channelWaiter
}

func NewChannel() *Channel {
	return &Channel{}
}

func (m *Channel) Lock(ctx context.Context) error {
	m.mu.Lock()
	if !m.locked && len(m.waiters) == 0 {
		m.locked = true
		m.mu.Unlock()
		return nil
	}

	waiter := &channelWaiter{ready: make(chan struct{}, 1)}
	m.waiters = append(m.waiters, waiter)
	m.mu.Unlock()

	select {
	case <-waiter.ready:
		return nil
	case <-ctx.Done():
		m.mu.Lock()
		var removed bool
		m.waiters, removed = removeChannelWaiter(m.waiters, waiter)
		m.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.ready
		return nil
	}
}

func (m *Channel) Unlock() {
	m.mu.Lock()
	if len(m.waiters) == 0 {
		m.locked = false
		m.mu.Unlock()
		return
	}

	waiter := m.waiters[0]
	m.waiters = m.waiters[1:]
	m.mu.Unlock()
	waiter.ready <- struct{}{}
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
