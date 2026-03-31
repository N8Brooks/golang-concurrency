package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Footman struct {
	footman *semaphore.Weighted
	forks   [5]*semaphore.Weighted
}

func NewFootman() *Footman {
	footman := semaphore.NewWeighted(4)
	forks := [5]*semaphore.Weighted{}
	for i := range 5 {
		forks[i] = semaphore.NewWeighted(1)
	}
	return &Footman{footman, forks}
}

func (dp *Footman) Dine(ctx context.Context, philosopher int, think, eat func()) {
	think()

	if err := dp.footman.Acquire(ctx, 1); err != nil {
		return
	}

	right := philosopher
	left := (philosopher + 1) % 5

	if err := dp.forks[right].Acquire(ctx, 1); err != nil {
		dp.footman.Release(1)
		return
	}
	if err := dp.forks[left].Acquire(ctx, 1); err != nil {
		dp.forks[right].Release(1)
		dp.footman.Release(1)
		return
	}

	eat()

	dp.forks[left].Release(1)
	dp.forks[right].Release(1)
	dp.footman.Release(1)
}
