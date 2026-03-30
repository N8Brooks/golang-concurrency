package solutions

import "sync"

type SyncCond struct {
	n       int
	mu      sync.Mutex
	cond    *sync.Cond
	fooTurn bool
}

func NewSyncCond(n int) *SyncCond {
	fb := &SyncCond{
		n:       n,
		fooTurn: true,
	}
	fb.cond = sync.NewCond(&fb.mu)
	return fb
}

func (fb *SyncCond) Foo(printFoo func()) {
	for range fb.n {
		fb.mu.Lock()
		for !fb.fooTurn {
			fb.cond.Wait()
		}
		printFoo()
		fb.fooTurn = false
		fb.cond.Broadcast()
		fb.mu.Unlock()
	}
}

func (fb *SyncCond) Bar(printBar func()) {
	for range fb.n {
		fb.mu.Lock()
		for fb.fooTurn {
			fb.cond.Wait()
		}
		printBar()
		fb.fooTurn = true
		fb.cond.Broadcast()
		fb.mu.Unlock()
	}
}
