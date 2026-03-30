package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Lightswitch struct {
	readers int
	mutex   *semaphore.Weighted
}

func NewLightswitch() *Lightswitch {
	return &Lightswitch{
		mutex: semaphore.NewWeighted(1),
	}
}

func (ls *Lightswitch) Lock(ctx context.Context, roomEmpty *semaphore.Weighted) error {
	if err := ls.mutex.Acquire(ctx, 1); err != nil {
		return err
	}

	ls.readers++
	if ls.readers == 1 {
		if err := roomEmpty.Acquire(ctx, 1); err != nil {
			ls.readers--
			ls.mutex.Release(1)
			return err
		}
	}

	ls.mutex.Release(1)
	return nil
}

func (ls *Lightswitch) Unlock(roomEmpty *semaphore.Weighted) {
	if err := ls.mutex.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}

	ls.readers--
	if ls.readers == 0 {
		roomEmpty.Release(1)
	}

	ls.mutex.Release(1)
}

type Semaphore struct {
	readSwitch *Lightswitch
	roomEmpty  *semaphore.Weighted
	turnstile  *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		readSwitch: NewLightswitch(),
		roomEmpty:  semaphore.NewWeighted(1),
		turnstile:  semaphore.NewWeighted(1),
	}
}

func (rw *Semaphore) Reader(ctx context.Context, read func()) error {
	if err := rw.turnstile.Acquire(ctx, 1); err != nil {
		return err
	}
	rw.turnstile.Release(1)

	if err := rw.readSwitch.Lock(ctx, rw.roomEmpty); err != nil {
		return err
	}
	read()
	rw.readSwitch.Unlock(rw.roomEmpty)
	return nil
}

func (rw *Semaphore) Writer(ctx context.Context, write func()) error {
	if err := rw.turnstile.Acquire(ctx, 1); err != nil {
		return err
	}
	if err := rw.roomEmpty.Acquire(ctx, 1); err != nil {
		rw.turnstile.Release(1)
		return err
	}

	write()

	rw.turnstile.Release(1)
	rw.roomEmpty.Release(1)
	return nil
}
