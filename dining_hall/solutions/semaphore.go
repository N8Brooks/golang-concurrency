// Package solutions contains implementations of the dining hall problem.
package solutions

import (
	"context"
	"sync"
)

type Semaphore struct {
	eating       int
	readyToLeave int
	waitingLeave chan struct{}
	mutex        *semaphore
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		mutex: newSemaphore(1),
	}
}

func (h *Semaphore) Student(ctx context.Context, dine, leave func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := h.mutex.Wait(ctx); err != nil {
		return err
	}
	h.eating++
	if h.eating == 2 && h.readyToLeave == 1 && h.waitingLeave != nil {
		waiter := h.waitingLeave
		h.waitingLeave = nil
		h.readyToLeave--
		close(waiter)
	}
	h.mutex.Signal()

	dine()

	if err := h.mutex.Wait(context.Background()); err != nil {
		panic(err)
	}
	h.eating--
	h.readyToLeave++

	switch {
	case h.eating == 1 && h.readyToLeave == 1:
		waiter := make(chan struct{})
		h.waitingLeave = waiter
		h.mutex.Signal()

		select {
		case <-waiter:
		case <-ctx.Done():
			if err := h.mutex.Wait(context.Background()); err != nil {
				panic(err)
			}
			if h.waitingLeave == waiter {
				h.waitingLeave = nil
				h.readyToLeave--
				h.mutex.Signal()
				return ctx.Err()
			}
			h.mutex.Signal()
			<-waiter
		}
	case h.eating == 0 && h.readyToLeave == 2:
		waiter := h.waitingLeave
		h.waitingLeave = nil
		h.readyToLeave -= 2
		h.mutex.Signal()
		if waiter != nil {
			close(waiter)
		}
	default:
		h.readyToLeave--
		h.mutex.Signal()
	}

	leave()
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
