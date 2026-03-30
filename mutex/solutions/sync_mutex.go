package solutions

import (
	"context"
	"sync"
)

type SyncMutex struct {
	mu     sync.Mutex
	cond   *sync.Cond
	locked bool
}

func NewSyncMutex() *SyncMutex {
	m := &SyncMutex{}
	m.cond = sync.NewCond(&m.mu)
	return m
}

func (m *SyncMutex) Lock(ctx context.Context) error {
	stop := context.AfterFunc(ctx, func() {
		m.mu.Lock()
		m.cond.Broadcast()
		m.mu.Unlock()
	})
	defer stop()

	m.mu.Lock()
	defer m.mu.Unlock()

	for m.locked {
		if err := ctx.Err(); err != nil {
			return err
		}
		m.cond.Wait()
	}

	m.locked = true
	return nil
}

func (m *SyncMutex) Unlock() {
	m.mu.Lock()
	m.locked = false
	m.cond.Signal()
	m.mu.Unlock()
}
