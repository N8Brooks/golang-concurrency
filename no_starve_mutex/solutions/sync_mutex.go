package solutions

import (
	"context"
	"sync"
)

type SyncMutex struct {
	mu       sync.Mutex
	cond     *sync.Cond
	next     uint64
	serving  uint64
	locked   bool
	canceled map[uint64]struct{}
}

func NewSyncMutex() *SyncMutex {
	m := &SyncMutex{
		canceled: make(map[uint64]struct{}),
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

func (m *SyncMutex) Lock(ctx context.Context) error {
	m.mu.Lock()
	ticket := m.next
	m.next++

	stop := context.AfterFunc(ctx, func() {
		m.mu.Lock()
		m.canceled[ticket] = struct{}{}
		m.advanceLocked()
		m.cond.Broadcast()
		m.mu.Unlock()
	})
	defer stop()

	for {
		m.advanceLocked()
		if !m.locked && ticket == m.serving {
			delete(m.canceled, ticket)
			m.locked = true
			m.mu.Unlock()
			return nil
		}
		if err := ctx.Err(); err != nil {
			m.canceled[ticket] = struct{}{}
			m.advanceLocked()
			m.cond.Broadcast()
			m.mu.Unlock()
			return err
		}
		m.cond.Wait()
	}
}

func (m *SyncMutex) Unlock() {
	m.mu.Lock()
	m.locked = false
	m.serving++
	m.advanceLocked()
	m.cond.Broadcast()
	m.mu.Unlock()
}

func (m *SyncMutex) advanceLocked() {
	for {
		if _, ok := m.canceled[m.serving]; !ok || m.locked {
			return
		}
		delete(m.canceled, m.serving)
		m.serving++
	}
}
