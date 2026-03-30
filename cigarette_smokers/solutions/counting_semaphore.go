package solutions

import (
	"context"
	"sync"
)

type countingSemaphore struct {
	mu      sync.Mutex
	count   int
	waiters []chan struct{}
}

func newCountingSemaphore(count int) *countingSemaphore {
	return &countingSemaphore{count: count}
}

func (s *countingSemaphore) Signal() {
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

func (s *countingSemaphore) Wait(ctx context.Context) bool {
	s.mu.Lock()
	if s.count > 0 {
		s.count--
		s.mu.Unlock()
		return true
	}
	if ctx.Err() != nil {
		s.mu.Unlock()
		return false
	}

	waiter := make(chan struct{})
	s.waiters = append(s.waiters, waiter)
	s.mu.Unlock()

	select {
	case <-waiter:
		return true
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
			return false
		}
		s.mu.Unlock()
		<-waiter
		return true
	}
}
