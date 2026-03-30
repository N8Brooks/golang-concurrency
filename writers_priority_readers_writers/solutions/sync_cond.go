package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	readers        int
	writer         bool
	waitingWriters int
	mu             sync.Mutex
	cond           *sync.Cond
}

func NewSyncCond() *SyncCond {
	rw := &SyncCond{}
	rw.cond = sync.NewCond(&rw.mu)
	return rw
}

func (rw *SyncCond) Reader(ctx context.Context, read func()) error {
	rw.mu.Lock()
	if err := rw.wait(ctx, func() bool {
		return !rw.writer && rw.waitingWriters == 0
	}); err != nil {
		rw.mu.Unlock()
		return err
	}
	rw.readers++
	rw.mu.Unlock()

	read()

	rw.mu.Lock()
	rw.readers--
	if rw.readers == 0 {
		rw.cond.Broadcast()
	}
	rw.mu.Unlock()
	return nil
}

func (rw *SyncCond) Writer(ctx context.Context, write func()) error {
	rw.mu.Lock()
	rw.waitingWriters++
	if err := rw.wait(ctx, func() bool {
		return !rw.writer && rw.readers == 0
	}); err != nil {
		rw.waitingWriters--
		rw.cond.Broadcast()
		rw.mu.Unlock()
		return err
	}
	rw.waitingWriters--
	rw.writer = true
	rw.mu.Unlock()

	write()

	rw.mu.Lock()
	rw.writer = false
	rw.cond.Broadcast()
	rw.mu.Unlock()
	return nil
}

func (rw *SyncCond) wait(ctx context.Context, ready func() bool) error {
	stop := context.AfterFunc(ctx, rw.cond.Broadcast)
	defer stop()

	for !ready() {
		if err := ctx.Err(); err != nil {
			return err
		}
		rw.cond.Wait()
	}

	return nil
}
