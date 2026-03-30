// Package solutions contains implementations of the multiplex problem.
package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	slots *semaphore.Weighted
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		slots: semaphore.NewWeighted(int64(limit)),
	}
}

func (m *Semaphore) Run(ctx context.Context, criticalSection func()) {
	if err := m.slots.Acquire(ctx, 1); err != nil {
		return
	}
	defer m.slots.Release(1)

	criticalSection()
}
