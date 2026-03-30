// Package solutions contains implementations of the sushi bar problem.
package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	capacity int
	eating   int
	waiting  int
	mustWait bool
	mu       sync.Mutex
	cond     *sync.Cond
}

func NewSyncCond(capacity int) *SyncCond {
	s := &SyncCond{capacity: capacity}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *SyncCond) Customer(ctx context.Context, dine func()) error {
	s.mu.Lock()
	if err := ctx.Err(); err != nil {
		s.mu.Unlock()
		return err
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			s.cond.Broadcast()
			s.mu.Unlock()
		case <-done:
		}
	}()

	for s.mustWait && ctx.Err() == nil {
		s.waiting++
		for s.mustWait && ctx.Err() == nil {
			s.cond.Wait()
		}
		s.waiting--
	}
	if err := ctx.Err(); err != nil {
		s.mu.Unlock()
		return err
	}

	s.eating++
	s.mustWait = s.eating == s.capacity
	s.mu.Unlock()

	dine()

	s.mu.Lock()
	s.eating--
	if s.eating == 0 {
		s.mustWait = false
		s.cond.Broadcast()
	}
	s.mu.Unlock()

	return nil
}
