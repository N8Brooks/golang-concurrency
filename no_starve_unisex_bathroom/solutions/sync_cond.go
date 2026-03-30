package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	mu            sync.Mutex
	cond          *sync.Cond
	turnstileHeld bool
	maleTotal     int
	femaleTotal   int
	maleInside    int
	femaleInside  int
}

func NewSyncCond() *SyncCond {
	b := &SyncCond{}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *SyncCond) Male(ctx context.Context, bathroom func()) error {
	if err := b.enterMale(ctx); err != nil {
		return err
	}
	bathroom()
	b.leaveMale()
	return nil
}

func (b *SyncCond) Female(ctx context.Context, bathroom func()) error {
	if err := b.enterFemale(ctx); err != nil {
		return err
	}
	bathroom()
	b.leaveFemale()
	return nil
}

func (b *SyncCond) enterMale(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	stop := context.AfterFunc(ctx, b.cond.Broadcast)
	defer stop()

	for b.turnstileHeld {
		if err := ctx.Err(); err != nil {
			return err
		}
		b.cond.Wait()
	}
	b.turnstileHeld = true

	for b.femaleTotal > 0 {
		if err := ctx.Err(); err != nil {
			b.turnstileHeld = false
			b.cond.Broadcast()
			return err
		}
		b.cond.Wait()
	}

	b.maleTotal++
	b.turnstileHeld = false
	b.cond.Broadcast()

	for b.maleInside == 3 {
		if err := ctx.Err(); err != nil {
			b.maleTotal--
			if b.maleTotal == 0 {
				b.cond.Broadcast()
			}
			b.cond.Broadcast()
			return err
		}
		b.cond.Wait()
	}

	b.maleInside++
	return nil
}

func (b *SyncCond) enterFemale(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	stop := context.AfterFunc(ctx, b.cond.Broadcast)
	defer stop()

	for b.turnstileHeld {
		if err := ctx.Err(); err != nil {
			return err
		}
		b.cond.Wait()
	}
	b.turnstileHeld = true

	for b.maleTotal > 0 {
		if err := ctx.Err(); err != nil {
			b.turnstileHeld = false
			b.cond.Broadcast()
			return err
		}
		b.cond.Wait()
	}

	b.femaleTotal++
	b.turnstileHeld = false
	b.cond.Broadcast()

	for b.femaleInside == 3 {
		if err := ctx.Err(); err != nil {
			b.femaleTotal--
			if b.femaleTotal == 0 {
				b.cond.Broadcast()
			}
			b.cond.Broadcast()
			return err
		}
		b.cond.Wait()
	}

	b.femaleInside++
	return nil
}

func (b *SyncCond) leaveMale() {
	b.mu.Lock()
	b.maleInside--
	b.maleTotal--
	b.cond.Broadcast()
	b.mu.Unlock()
}

func (b *SyncCond) leaveFemale() {
	b.mu.Lock()
	b.femaleInside--
	b.femaleTotal--
	b.cond.Broadcast()
	b.mu.Unlock()
}
