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
	mu      sync.Mutex
	locked  bool
	waiters []*semaphoreWaiter
}

func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

func (m *Semaphore) Lock(ctx context.Context) error {
	m.mu.Lock()
	if !m.locked && len(m.waiters) == 0 {
		m.locked = true
		m.mu.Unlock()
		return nil
	}

	waiter := newSemaphoreWaiter()
	m.waiters = append(m.waiters, waiter)
	m.mu.Unlock()

	if err := waiter.ready.Acquire(ctx, 1); err != nil {
		m.mu.Lock()
		var removed bool
		m.waiters, removed = removeSemaphoreWaiter(m.waiters, waiter)
		m.mu.Unlock()
		if removed {
			return err
		}
		if err := waiter.ready.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
	}
	return nil
}

func (m *Semaphore) Unlock() {
	m.mu.Lock()
	if len(m.waiters) == 0 {
		m.locked = false
		m.mu.Unlock()
		return
	}

	waiter := m.waiters[0]
	m.waiters = m.waiters[1:]
	m.mu.Unlock()
	waiter.ready.Release(1)
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
