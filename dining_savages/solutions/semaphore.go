// Package solutions contains implementations of the dining savages problem.
package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	capacity int
	servings int
	mutex    *semaphore.Weighted
	emptyPot *semaphore.Weighted
	fullPot  *semaphore.Weighted
}

func NewSemaphore(capacity int) *Semaphore {
	ds := &Semaphore{
		capacity: capacity,
		mutex:    semaphore.NewWeighted(1),
		emptyPot: semaphore.NewWeighted(1),
		fullPot:  semaphore.NewWeighted(1),
	}
	ctx := context.Background()
	if err := ds.emptyPot.Acquire(ctx, 1); err != nil {
		panic(err)
	}
	if err := ds.fullPot.Acquire(ctx, 1); err != nil {
		panic(err)
	}
	return ds
}

func (ds *Semaphore) Savage(ctx context.Context, getServingFromPot, eat func()) {
	for {
		if err := ds.mutex.Acquire(ctx, 1); err != nil {
			return
		}
		if ds.servings == 0 {
			ds.emptyPot.Release(1)
			if err := ds.fullPot.Acquire(ctx, 1); err != nil {
				ds.mutex.Release(1)
				return
			}
			ds.servings = ds.capacity
		}
		ds.servings--
		getServingFromPot()
		ds.mutex.Release(1)

		eat()
	}
}

func (ds *Semaphore) Cook(ctx context.Context, putServingsInPot func(servings int)) {
	for {
		if err := ds.emptyPot.Acquire(ctx, 1); err != nil {
			return
		}
		putServingsInPot(ds.capacity)
		ds.fullPot.Release(1)
	}
}
