package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	token *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		token: semaphore.NewWeighted(1),
	}
}

func (m *Semaphore) Lock(ctx context.Context) error {
	return m.token.Acquire(ctx, 1)
}

func (m *Semaphore) Unlock() {
	m.token.Release(1)
}
