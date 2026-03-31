package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Asymmetric struct {
	forks [5]*semaphore.Weighted
}

func NewAsymmetric() *Asymmetric {
	forks := [5]*semaphore.Weighted{}
	for i := range 5 {
		forks[i] = semaphore.NewWeighted(1)
	}
	return &Asymmetric{forks}
}

func (dp *Asymmetric) Dine(ctx context.Context, philosopher int, think, eat func()) {
	think()

	first := philosopher
	second := (philosopher + 1) % 5
	if philosopher == 0 {
		first, second = second, first
	}

	if err := dp.forks[first].Acquire(ctx, 1); err != nil {
		return
	}
	if err := dp.forks[second].Acquire(ctx, 1); err != nil {
		dp.forks[first].Release(1)
		return
	}

	eat()

	dp.forks[second].Release(1)
	dp.forks[first].Release(1)
}
