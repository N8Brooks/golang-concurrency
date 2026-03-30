//go:build challenge

// Package sushibar contains the challenge version of the sushi bar problem.
//
// A fixed number of seats are available. Customers may sit immediately while
// seats remain open, but once the bar is full, later arrivals must wait for the
// entire current party to leave before sitting down.
package sushibar

import (
	"context"
	"sync"
)

type SushiBar struct {
	capacity int
	eating   int
	waiting  int
	mustWait bool
	mutex    *semaphore
	block    *semaphore
}

func NewSushiBar(capacity int) *SushiBar {
	return &SushiBar{
		capacity: capacity,
		mutex:    newSemaphore(1),
		block:    newSemaphore(0),
	}
}

func (s *SushiBar) Customer(ctx context.Context, dine func()) error {
	if err := s.mutex.Wait(ctx); err != nil {
		return err
	}

	if s.mustWait {
		s.waiting++
		s.mutex.Signal()

		if err := s.block.Wait(ctx); err != nil {
			return err
		}
	} else {
		s.eating++
		s.mustWait = s.eating == s.capacity
		s.mutex.Signal()
	}

	dine()

	if err := s.mutex.Wait(context.Background()); err != nil {
		panic(err)
	}
	s.eating--
	if s.eating == 0 {
		n := min(s.capacity, s.waiting)
		s.waiting -= n
		s.eating += n
		s.mustWait = s.eating == s.capacity
		for range n {
			s.block.Signal()
		}
	}
	s.mutex.Signal()

	return nil
}

type semaphore struct {
	mu      sync.Mutex
	count   int
	waiters []chan struct{}
}

func newSemaphore(count int) *semaphore {
	return &semaphore{count: count}
}

func (s *semaphore) Signal() {
	s.mu.Lock()
	if len(s.waiters) > 0 {
		waiter := s.waiters[0]
		copy(s.waiters, s.waiters[1:])
		s.waiters[len(s.waiters)-1] = nil
		s.waiters = s.waiters[:len(s.waiters)-1]
		s.mu.Unlock()
		close(waiter)
		return
	}
	s.count++
	s.mu.Unlock()
}

func (s *semaphore) Wait(ctx context.Context) error {
	s.mu.Lock()
	if s.count > 0 {
		s.count--
		s.mu.Unlock()
		return nil
	}
	if err := ctx.Err(); err != nil {
		s.mu.Unlock()
		return err
	}

	waiter := make(chan struct{})
	s.waiters = append(s.waiters, waiter)
	s.mu.Unlock()

	select {
	case <-waiter:
		return nil
	case <-ctx.Done():
		s.mu.Lock()
		index := -1
		for i, candidate := range s.waiters {
			if candidate == waiter {
				index = i
				break
			}
		}
		if index >= 0 {
			copy(s.waiters[index:], s.waiters[index+1:])
			s.waiters[len(s.waiters)-1] = nil
			s.waiters = s.waiters[:len(s.waiters)-1]
			s.mu.Unlock()
			return ctx.Err()
		}
		s.mu.Unlock()
		<-waiter
		return nil
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
