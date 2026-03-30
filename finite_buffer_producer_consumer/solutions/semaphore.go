package solutions

import (
	"context"
	"sync"
)

type Semaphore struct {
	mu     *semaphore
	items  *semaphore
	spaces *semaphore
	buffer []int
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		mu:     newSemaphore(1),
		items:  newSemaphore(0),
		spaces: newSemaphore(capacity),
	}
}

func (pc *Semaphore) Produce(ctx context.Context, waitForEvent func() int) error {
	event := waitForEvent()

	if err := pc.spaces.Wait(ctx); err != nil {
		return err
	}
	if err := pc.mu.Wait(context.Background()); err != nil {
		panic(err)
	}
	pc.buffer = append(pc.buffer, event)
	pc.mu.Signal()

	pc.items.Signal()
	return nil
}

func (pc *Semaphore) Consume(ctx context.Context, process func(int)) error {
	if err := pc.items.Wait(ctx); err != nil {
		return err
	}
	if err := pc.mu.Wait(context.Background()); err != nil {
		panic(err)
	}
	event := pc.buffer[0]
	pc.buffer = pc.buffer[1:]
	pc.mu.Signal()

	pc.spaces.Signal()
	process(event)
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
