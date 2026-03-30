package solutions

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	n   int
	foo *semaphore.Weighted
	bar *semaphore.Weighted
}

func NewSemaphore(n int) *Semaphore {
	bar := semaphore.NewWeighted(1)
	if err := bar.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	return &Semaphore{
		n:   n,
		foo: semaphore.NewWeighted(1),
		bar: bar,
	}
}

func (fb *Semaphore) Foo(printFoo func()) {
	for range fb.n {
		if err := fb.foo.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
		printFoo()
		fb.bar.Release(1)
	}
}

func (fb *Semaphore) Bar(printBar func()) {
	for range fb.n {
		if err := fb.bar.Acquire(context.Background(), 1); err != nil {
			panic(err)
		}
		printBar()
		fb.foo.Release(1)
	}
}
