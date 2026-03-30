package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type semaphoreBarrier struct {
	n          int64
	count      int64
	mutex      *semaphore.Weighted
	turnstile1 *semaphore.Weighted
	turnstile2 *semaphore.Weighted
}

func newSemaphoreBarrier(n int64) *semaphoreBarrier {
	turnstile1 := semaphore.NewWeighted(n)
	turnstile2 := semaphore.NewWeighted(n)
	if err := turnstile1.Acquire(context.Background(), n); err != nil {
		panic(err)
	}
	if err := turnstile2.Acquire(context.Background(), n); err != nil {
		panic(err)
	}
	return &semaphoreBarrier{
		n:          n,
		mutex:      semaphore.NewWeighted(1),
		turnstile1: turnstile1,
		turnstile2: turnstile2,
	}
}

func (b *semaphoreBarrier) Wait() {
	if err := b.mutex.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	b.count++
	if b.count == b.n {
		b.turnstile1.Release(b.n)
	}
	b.mutex.Release(1)

	if err := b.turnstile1.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}

	if err := b.mutex.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	b.count--
	if b.count == 0 {
		b.turnstile2.Release(b.n)
	}
	b.mutex.Release(1)

	if err := b.turnstile2.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
}

type Semaphore struct {
	mutex      *semaphore.Weighted
	oxygen     int64
	hydrogen   int64
	barrier    *semaphoreBarrier
	oxyQueue   *semaphore.Weighted
	hydroQueue *semaphore.Weighted
}

func NewSemaphore() *Semaphore {
	oxyQueue := semaphore.NewWeighted(1)
	hydroQueue := semaphore.NewWeighted(2)
	if err := oxyQueue.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	if err := hydroQueue.Acquire(context.Background(), 2); err != nil {
		panic(err)
	}
	return &Semaphore{
		mutex:      semaphore.NewWeighted(1),
		barrier:    newSemaphoreBarrier(3),
		oxyQueue:   oxyQueue,
		hydroQueue: hydroQueue,
	}
}

func (h2o *Semaphore) Hydrogen(bond func()) {
	if err := h2o.mutex.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	h2o.hydrogen++
	if h2o.hydrogen >= 2 && h2o.oxygen >= 1 {
		h2o.hydroQueue.Release(2)
		h2o.hydrogen -= 2
		h2o.oxyQueue.Release(1)
		h2o.oxygen--
	} else {
		h2o.mutex.Release(1)
	}

	if err := h2o.hydroQueue.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	bond()
	h2o.barrier.Wait()
}

func (h2o *Semaphore) Oxygen(bond func()) {
	if err := h2o.mutex.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	h2o.oxygen++
	if h2o.hydrogen >= 2 {
		h2o.hydroQueue.Release(2)
		h2o.hydrogen -= 2
		h2o.oxyQueue.Release(1)
		h2o.oxygen--
	} else {
		h2o.mutex.Release(1)
	}

	if err := h2o.oxyQueue.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	bond()
	h2o.barrier.Wait()
	h2o.mutex.Release(1)
}
