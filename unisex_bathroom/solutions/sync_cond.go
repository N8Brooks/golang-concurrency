package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	mu       sync.Mutex
	cond     *sync.Cond
	inside   int
	occupant int
}

const (
	none = iota
	male
	female
)

func NewSyncCond() *SyncCond {
	b := &SyncCond{}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *SyncCond) Male(ctx context.Context, bathroom func()) error {
	if err := b.enter(ctx, male); err != nil {
		return err
	}
	bathroom()
	b.leave(male)
	return nil
}

func (b *SyncCond) Female(ctx context.Context, bathroom func()) error {
	if err := b.enter(ctx, female); err != nil {
		return err
	}
	bathroom()
	b.leave(female)
	return nil
}

func (b *SyncCond) enter(ctx context.Context, kind int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	stop := context.AfterFunc(ctx, b.cond.Broadcast)
	defer stop()

	for b.inside == 3 || (b.occupant != none && b.occupant != kind) {
		if err := ctx.Err(); err != nil {
			return err
		}
		b.cond.Wait()
	}

	b.inside++
	b.occupant = kind
	return nil
}

func (b *SyncCond) leave(kind int) {
	b.mu.Lock()
	b.inside--
	if b.inside == 0 {
		b.occupant = none
	}
	b.cond.Broadcast()
	b.mu.Unlock()
}
