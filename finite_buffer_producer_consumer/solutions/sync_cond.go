package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	capacity int
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	buffer   []int
}

func NewSyncCond(capacity int) *SyncCond {
	pc := &SyncCond{
		capacity: capacity,
	}
	pc.notEmpty = sync.NewCond(&pc.mu)
	pc.notFull = sync.NewCond(&pc.mu)
	return pc
}

func (pc *SyncCond) Produce(ctx context.Context, waitForEvent func() int) error {
	event := waitForEvent()

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
			pc.notFull.Broadcast()
			pc.mu.Unlock()
		case <-done:
		}
	}()

	for len(pc.buffer) == pc.capacity && ctx.Err() == nil {
		pc.notFull.Wait()
	}
	if err := ctx.Err(); err != nil {
		pc.mu.Unlock()
		return err
	}

	pc.buffer = append(pc.buffer, event)
	pc.mu.Unlock()
	pc.notEmpty.Signal()
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
			pc.notEmpty.Broadcast()
			pc.mu.Unlock()
		case <-done:
		}
	}()

	for len(pc.buffer) == 0 && ctx.Err() == nil {
		pc.notEmpty.Wait()
	}
	if err := ctx.Err(); err != nil {
		pc.mu.Unlock()
		return err
	}

	event := pc.buffer[0]
	pc.buffer = pc.buffer[1:]
	pc.mu.Unlock()
	pc.notFull.Signal()

	process(event)
	return nil
}
