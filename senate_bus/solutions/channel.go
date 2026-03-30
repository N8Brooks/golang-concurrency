package solutions

import (
	"context"
	"sync"
)

type riderWaiter struct {
	ready   chan struct{}
	boarded chan struct{}
}

type Channel struct {
	capacity int
	mu       sync.Mutex
	waiting  []*riderWaiter
}

func NewChannel(capacity int) *Channel {
	return &Channel{capacity: capacity}
}

func (b *Channel) Rider(ctx context.Context, boardBus func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	waiter := &riderWaiter{
		ready:   make(chan struct{}),
		boarded: make(chan struct{}),
	}

	b.mu.Lock()
	b.waiting = append(b.waiting, waiter)
	b.mu.Unlock()

	select {
	case <-waiter.ready:
		boardBus()
		close(waiter.boarded)
		return nil
	case <-ctx.Done():
		b.mu.Lock()
		removed := removeWaiter(&b.waiting, waiter)
		b.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.ready
		boardBus()
		close(waiter.boarded)
		return nil
	}
}

func (b *Channel) Bus(ctx context.Context, depart func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.mu.Lock()
	n := min(len(b.waiting), b.capacity)
	if n == 0 {
		b.mu.Unlock()
		depart()
		return nil
	}
	selected := append([]*riderWaiter(nil), b.waiting[:n]...)
	b.waiting = b.waiting[n:]
	b.mu.Unlock()

	for _, waiter := range selected {
		close(waiter.ready)
	}
	for _, waiter := range selected {
		<-waiter.boarded
	}

	depart()
	return nil
}

func removeWaiter(waiters *[]*riderWaiter, target *riderWaiter) bool {
	for i, waiter := range *waiters {
		if waiter != target {
			continue
		}
		copy((*waiters)[i:], (*waiters)[i+1:])
		last := len(*waiters) - 1
		(*waiters)[last] = nil
		*waiters = (*waiters)[:last]
		return true
	}
	return false
}
