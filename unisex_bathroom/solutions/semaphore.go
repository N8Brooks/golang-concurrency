package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Lightswitch struct {
	count int64
	mutex *semaphore.Weighted
}

func NewLightswitch() *Lightswitch {
	return &Lightswitch{
		mutex: semaphore.NewWeighted(1),
	}
}

func (ls *Lightswitch) Lock(ctx context.Context, sem *semaphore.Weighted) error {
	if err := ls.mutex.Acquire(ctx, 1); err != nil {
		return err
	}

	ls.count++
	if ls.count == 1 {
		if err := sem.Acquire(ctx, 1); err != nil {
			ls.count--
			ls.mutex.Release(1)
			return err
		}
	}

	ls.mutex.Release(1)
	return nil
}

func (ls *Lightswitch) Unlock(sem *semaphore.Weighted) {
	if err := ls.mutex.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}

	ls.count--
	if ls.count == 0 {
		sem.Release(1)
	}

	ls.mutex.Release(1)
}

type Semaphore struct {
	empty           *semaphore.Weighted
	maleSwitch      *Lightswitch
	femaleSwitch    *Lightswitch
	maleMultiplex   *semaphore.Weighted
	femaleMultiplex *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		empty:           semaphore.NewWeighted(1),
		maleSwitch:      NewLightswitch(),
		femaleSwitch:    NewLightswitch(),
		maleMultiplex:   semaphore.NewWeighted(3),
		femaleMultiplex: semaphore.NewWeighted(3),
	}
}

func (b *Semaphore) Male(ctx context.Context, bathroom func()) error {
	if err := b.maleSwitch.Lock(ctx, b.empty); err != nil {
		return err
	}
	if err := b.maleMultiplex.Acquire(ctx, 1); err != nil {
		b.maleSwitch.Unlock(b.empty)
		return err
	}

	bathroom()

	b.maleMultiplex.Release(1)
	b.maleSwitch.Unlock(b.empty)
	return nil
}

func (b *Semaphore) Female(ctx context.Context, bathroom func()) error {
	if err := b.femaleSwitch.Lock(ctx, b.empty); err != nil {
		return err
	}
	if err := b.femaleMultiplex.Acquire(ctx, 1); err != nil {
		b.femaleSwitch.Unlock(b.empty)
		return err
	}

	bathroom()

	b.femaleMultiplex.Release(1)
	b.femaleSwitch.Unlock(b.empty)
	return nil
}
