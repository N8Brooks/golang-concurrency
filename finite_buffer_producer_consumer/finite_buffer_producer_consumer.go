//go:build challenge

// Package finitebufferproducerconsumer contains the challenge version of the
// finite-buffer producer-consumer problem.
//
// Producers wait for events and append them to a bounded shared buffer.
// Consumers remove buffered events and process them. Producers must wait while
// the buffer is full, consumers must wait while it is empty, and neither
// waiting nor processing should happen while holding the buffer lock.
package finitebufferproducerconsumer

import (
	"context"
	"sync"
)

type ProducerConsumer struct {
	mu     sync.Mutex
	items  *countingSemaphore
	spaces *countingSemaphore
	buffer []int
}

func NewProducerConsumer(capacity int) *ProducerConsumer {
	return &ProducerConsumer{
		items:  newCountingSemaphore(0),
		spaces: newCountingSemaphore(capacity),
	}
}

func (pc *ProducerConsumer) Produce(ctx context.Context, waitForEvent func() int) error {
	event := waitForEvent()

	if err := pc.spaces.Wait(ctx); err != nil {
		return err
	}

	pc.mu.Lock()
	pc.buffer = append(pc.buffer, event)
	pc.mu.Unlock()

	pc.items.Signal()
	return nil
}

func (pc *ProducerConsumer) Consume(ctx context.Context, process func(int)) error {
	if err := pc.items.Wait(ctx); err != nil {
		return err
	}

	pc.mu.Lock()
	event := pc.buffer[0]
	pc.buffer = pc.buffer[1:]
	pc.mu.Unlock()

	pc.spaces.Signal()
	process(event)
	return nil
}

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

func (s *countingSemaphore) Wait(ctx context.Context) error {
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
