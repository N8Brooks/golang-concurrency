// Package solutions contains implementations of the dining savages problem.
package solutions

import (
	"context"
	"sync"
)

type Channel struct {
	capacity int
	servings int
	mu       sync.Mutex
	emptyPot chan struct{}
	fullPot  chan struct{}
}

func NewChannel(capacity int) *Channel {
	return &Channel{
		capacity: capacity,
		emptyPot: make(chan struct{}),
		fullPot:  make(chan struct{}),
	}
}

func (ds *Channel) Savage(ctx context.Context, getServingFromPot, eat func()) {
	for {
		if ctx.Err() != nil {
			return
		}

		ds.mu.Lock()
		if ds.servings == 0 {
			select {
			case <-ctx.Done():
				ds.mu.Unlock()
				return
			case ds.emptyPot <- struct{}{}:
			}
			select {
			case <-ctx.Done():
				ds.mu.Unlock()
				return
			case <-ds.fullPot:
			}
			ds.servings = ds.capacity
		}
		ds.servings--
		getServingFromPot()
		ds.mu.Unlock()

		eat()
	}
}

func (ds *Channel) Cook(ctx context.Context, putServingsInPot func(servings int)) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ds.emptyPot:
		}

		putServingsInPot(ds.capacity)

		select {
		case <-ctx.Done():
			return
		case ds.fullPot <- struct{}{}:
		}
	}
}
