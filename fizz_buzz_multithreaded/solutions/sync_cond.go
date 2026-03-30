// Package solutions contains implementations of the multithreaded FizzBuzz
// problem.
package solutions

import "sync"

type SyncCond struct {
	n       int
	current int
	mu      sync.Mutex
	cond    *sync.Cond
}

func NewSyncCond(n int) *SyncCond {
	fb := &SyncCond{
		n:       n,
		current: 1,
	}
	fb.cond = sync.NewCond(&fb.mu)
	return fb
}

func (fb *SyncCond) Fizz(printFizz func()) {
	fb.run(func(i int) bool {
		return i%3 == 0 && i%5 != 0
	}, func(int) {
		printFizz()
	})
}

func (fb *SyncCond) Buzz(printBuzz func()) {
	fb.run(func(i int) bool {
		return i%5 == 0 && i%3 != 0
	}, func(int) {
		printBuzz()
	})
}

func (fb *SyncCond) FizzBuzz(printFizzBuzz func()) {
	fb.run(func(i int) bool {
		return i%15 == 0
	}, func(int) {
		printFizzBuzz()
	})
}

func (fb *SyncCond) Number(printNumber func(int)) {
	fb.run(func(i int) bool {
		return i%3 != 0 && i%5 != 0
	}, printNumber)
}

func (fb *SyncCond) run(canPrint func(int) bool, print func(int)) {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	for fb.current <= fb.n {
		for fb.current <= fb.n && !canPrint(fb.current) {
			fb.cond.Wait()
		}
		if fb.current > fb.n {
			break
		}

		print(fb.current)
		fb.current++
		fb.cond.Broadcast()
	}

	fb.cond.Broadcast()
}
