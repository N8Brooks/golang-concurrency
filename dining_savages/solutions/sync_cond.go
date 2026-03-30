// Package solutions contains implementations of the dining savages problem.
package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	capacity int
	servings int
	refill   bool
	mu       sync.Mutex
	emptyPot *sync.Cond
	fullPot  *sync.Cond
}

func NewSyncCond(capacity int) *SyncCond {
	ds := &SyncCond{
		capacity: capacity,
	}
	ds.emptyPot = sync.NewCond(&ds.mu)
	ds.fullPot = sync.NewCond(&ds.mu)
	return ds
}

func (ds *SyncCond) Savage(ctx context.Context, getServingFromPot, eat func()) {
	for {
		ds.mu.Lock()
		if err := ctx.Err(); err != nil {
			ds.mu.Unlock()
			return
		}

		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				ds.mu.Lock()
				ds.emptyPot.Broadcast()
				ds.fullPot.Broadcast()
				ds.mu.Unlock()
			case <-done:
			}
		}()

		for ds.servings == 0 && ctx.Err() == nil {
			if !ds.refill {
				ds.refill = true
				ds.emptyPot.Signal()
			}
			ds.fullPot.Wait()
		}
		close(done)

		if err := ctx.Err(); err != nil {
			ds.mu.Unlock()
			return
		}

		ds.servings--
		getServingFromPot()
		ds.mu.Unlock()

		eat()
	}
}

func (ds *SyncCond) Cook(ctx context.Context, putServingsInPot func(servings int)) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			ds.mu.Lock()
			ds.emptyPot.Broadcast()
			ds.fullPot.Broadcast()
			ds.mu.Unlock()
		case <-done:
		}
	}()

	for {
		for !ds.refill && ctx.Err() == nil {
			ds.emptyPot.Wait()
		}
		if ctx.Err() != nil {
			return
		}

		putServingsInPot(ds.capacity)
		ds.servings = ds.capacity
		ds.refill = false
		ds.fullPot.Signal()
	}
}
