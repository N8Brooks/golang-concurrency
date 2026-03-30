package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Lightswitch struct {
	count int
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
	readSwitch  *Lightswitch
	writeSwitch *Lightswitch
	noReaders   *semaphore.Weighted
	noWriters   *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		readSwitch:  NewLightswitch(),
		writeSwitch: NewLightswitch(),
		noReaders:   semaphore.NewWeighted(1),
		noWriters:   semaphore.NewWeighted(1),
	}
}

func (rw *Semaphore) Reader(ctx context.Context, read func()) error {
	if err := rw.noReaders.Acquire(ctx, 1); err != nil {
		return err
	}
	if err := rw.readSwitch.Lock(ctx, rw.noWriters); err != nil {
		rw.noReaders.Release(1)
		return err
	}
	rw.noReaders.Release(1)

	read()

	rw.readSwitch.Unlock(rw.noWriters)
	return nil
}

func (rw *Semaphore) Writer(ctx context.Context, write func()) error {
	if err := rw.writeSwitch.Lock(ctx, rw.noReaders); err != nil {
		return err
	}
	if err := rw.noWriters.Acquire(ctx, 1); err != nil {
		rw.writeSwitch.Unlock(rw.noReaders)
		return err
	}

	write()

	rw.noWriters.Release(1)
	rw.writeSwitch.Unlock(rw.noReaders)
	return nil
}
