package solutions

import (
	"context"
	"math"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	mu     sync.Mutex
	buffer []int
	items  *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	items := semaphore.NewWeighted(math.MaxInt64)
	items.Acquire(context.Background(), math.MaxInt64)
	return &Semaphore{items: items}
}

func (pc *Semaphore) Produce(_ context.Context, waitForEvent func() int) {
	event := waitForEvent()
	pc.mu.Lock()
	pc.buffer = append(pc.buffer, event)
	pc.mu.Unlock()
	pc.items.Release(1)
}

func (pc *Semaphore) Consume(ctx context.Context, process func(int)) {
	if err := pc.items.Acquire(ctx, 1); err != nil {
		return
	}
	pc.mu.Lock()
	n := len(pc.buffer) - 1
	event := pc.buffer[n]
	pc.buffer = pc.buffer[:n]
	pc.mu.Unlock()
	process(event)
}
