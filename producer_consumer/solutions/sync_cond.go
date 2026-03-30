package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	mu     sync.Mutex
	cond   *sync.Cond
	buffer []int
}

func NewSyncCond() *SyncCond {
	pc := &SyncCond{}
	pc.cond = sync.NewCond(&pc.mu)
	return pc
}

func (pc *SyncCond) Produce(_ context.Context, waitForEvent func() int) error {
	event := waitForEvent()

	pc.mu.Lock()
	pc.buffer = append(pc.buffer, event)
	pc.mu.Unlock()
	pc.cond.Signal()
	return nil
}

func (pc *SyncCond) Consume(ctx context.Context, process func(int)) error {
	pc.mu.Lock()
	if err := ctx.Err(); err != nil {
		pc.mu.Unlock()
		return err
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			pc.mu.Lock()
			pc.cond.Broadcast()
			pc.mu.Unlock()
		case <-done:
		}
	}()

	for len(pc.buffer) == 0 && ctx.Err() == nil {
		pc.cond.Wait()
	}
	if err := ctx.Err(); err != nil {
		pc.mu.Unlock()
		return err
	}

	event := pc.buffer[0]
	pc.buffer = pc.buffer[1:]
	pc.mu.Unlock()

	process(event)
	return nil
}
