package diningsavages

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type DiningSavages struct {
	servings int
	mutex    *semaphore.Weighted
	emptyPot *semaphore.Weighted
	fullPot  *semaphore.Weighted
}

func NewDiningSavages() *DiningSavages {
	ds := &DiningSavages{
		mutex:    semaphore.NewWeighted(1), // Mutex to protect access to the pot
		emptyPot: semaphore.NewWeighted(1), // Semaphore to signal the cook to fill the pot
		fullPot:  semaphore.NewWeighted(1), // Semaphore to signal that the pot is full
	}
	ctx := context.Background()
	ds.emptyPot.Acquire(ctx, 1) // Start with the empty pot semaphore acquired
	ds.fullPot.Acquire(ctx, 1)  // Start with the full pot semaphore acquired
	return ds
}

func (ds *DiningSavages) Savage(ctx context.Context, getServingFromPot, eat func()) {
	for {
		if ds.mutex.Acquire(ctx, 1) != nil {
			return
		}
		if ds.servings == 0 {
			ds.emptyPot.Release(1) // Signal the cook to fill the pot
			if ds.fullPot.Acquire(ctx, 1) != nil {
				ds.mutex.Release(1) // Release the mutex before returning
				return
			}
		}
		ds.servings--
		getServingFromPot() // Get a serving from the pot
		ds.mutex.Release(1) // Release the mutex
		eat()               // Eat the serving
	}
}

func (ds *DiningSavages) Cook(ctx context.Context, fillPot func() int) {
	for {
		if ds.emptyPot.Acquire(ctx, 1) != nil {
			return
		}
		ds.servings = fillPot() // Fill the pot with new servings
		ds.fullPot.Release(1)   // Signal that the pot is full
	}
}
