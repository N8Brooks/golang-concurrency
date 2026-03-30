package solutions

import (
	"context"
	"sync"
)

type condWaiter struct {
	selected bool
	boarded  bool
}

type SyncCond struct {
	capacity int
	mu       sync.Mutex
	cond     *sync.Cond
	waiting  []*condWaiter
}

func NewSyncCond(capacity int) *SyncCond {
	b := &SyncCond{capacity: capacity}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *SyncCond) Rider(ctx context.Context, boardBus func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	waiter := &condWaiter{}

	b.mu.Lock()
	b.waiting = append(b.waiting, waiter)

	stop := context.AfterFunc(ctx, b.cond.Broadcast)
	defer stop()

	for !waiter.selected {
		if err := ctx.Err(); err != nil {
			if removeCondWaiter(&b.waiting, waiter) {
				b.mu.Unlock()
				return err
			}
		}
		b.cond.Wait()
	}
	b.mu.Unlock()

	boardBus()

	b.mu.Lock()
	waiter.boarded = true
	b.cond.Broadcast()
	b.mu.Unlock()
	return nil
}

func (b *SyncCond) Bus(ctx context.Context, depart func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.mu.Lock()
	n := min(len(b.waiting), b.capacity)
	if n == 0 {
		b.mu.Unlock()
		depart()
		return nil
	}

	selected := append([]*condWaiter(nil), b.waiting[:n]...)
	b.waiting = b.waiting[n:]
	for _, waiter := range selected {
		waiter.selected = true
	}
	b.cond.Broadcast()

	for !allBoarded(selected) {
		b.cond.Wait()
	}
	b.mu.Unlock()

	depart()
	return nil
}

func allBoarded(waiters []*condWaiter) bool {
	for _, waiter := range waiters {
		if !waiter.boarded {
			return false
		}
	}
	return true
}

func removeCondWaiter(waiters *[]*condWaiter, target *condWaiter) bool {
	for i, waiter := range *waiters {
		if waiter != target {
			continue
		}
		copy((*waiters)[i:], (*waiters)[i+1:])
		last := len(*waiters) - 1
		(*waiters)[last] = nil
		*waiters = (*waiters)[:last]
		return true
	}
	return false
}
