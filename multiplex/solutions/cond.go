// Package solutions contains implementations of the multiplex problem.
package solutions

import (
	"context"
	"sync"
)

type Cond struct {
	limit  int
	active int
	mu     sync.Mutex
	cond   *sync.Cond
}

func NewCond(limit int) *Cond {
	m := &Cond{
		limit: limit,
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

func (m *Cond) Run(ctx context.Context, criticalSection func()) {
	m.mu.Lock()
	if ctx.Err() != nil {
		m.mu.Unlock()
		return
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			m.mu.Lock()
			m.cond.Broadcast()
			m.mu.Unlock()
		case <-done:
		}
	}()

	for m.active >= m.limit && ctx.Err() == nil {
		m.cond.Wait()
	}
	if ctx.Err() != nil {
		m.mu.Unlock()
		return
	}

	m.active++
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.active--
		m.mu.Unlock()
		m.cond.Signal()
	}()

	criticalSection()
}
