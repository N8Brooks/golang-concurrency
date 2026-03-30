package solutions

import "sync"

type countingSemaphore struct {
	mu      sync.Mutex
	cond    *sync.Cond
	count   int
	waiters int
}

func newCountingSemaphore(count int) *countingSemaphore {
	s := &countingSemaphore{count: count}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *countingSemaphore) Signal(n int) {
	s.mu.Lock()
	s.count += n
	s.cond.Broadcast()
	s.mu.Unlock()
}

func (s *countingSemaphore) Wait() {
	s.mu.Lock()
	s.waiters++
	for s.count == 0 {
		s.cond.Wait()
	}
	s.waiters--
	s.count--
	s.mu.Unlock()
}

type condBarrier struct {
	n          int
	count      int
	generation int
	mu         sync.Mutex
	cond       *sync.Cond
}

func newCondBarrier(n int) *condBarrier {
	b := &condBarrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *condBarrier) Wait() {
	b.mu.Lock()
	generation := b.generation
	b.count++
	if b.count == b.n {
		b.count = 0
		b.generation++
		b.cond.Broadcast()
		b.mu.Unlock()
		return
	}
	for generation == b.generation {
		b.cond.Wait()
	}
	b.mu.Unlock()
}

type SyncCond struct {
	mu         sync.Mutex
	oxygen     int
	hydrogen   int
	barrier    *condBarrier
	oxyQueue   *countingSemaphore
	hydroQueue *countingSemaphore
}

func NewSyncCond() *SyncCond {
	return &SyncCond{
		barrier:    newCondBarrier(3),
		oxyQueue:   newCountingSemaphore(0),
		hydroQueue: newCountingSemaphore(0),
	}
}

func (h2o *SyncCond) Hydrogen(bond func()) {
	h2o.mu.Lock()
	h2o.hydrogen++
	if h2o.hydrogen >= 2 && h2o.oxygen >= 1 {
		h2o.hydroQueue.Signal(2)
		h2o.hydrogen -= 2
		h2o.oxyQueue.Signal(1)
		h2o.oxygen--
	} else {
		h2o.mu.Unlock()
	}

	h2o.hydroQueue.Wait()
	bond()
	h2o.barrier.Wait()
}

func (h2o *SyncCond) Oxygen(bond func()) {
	h2o.mu.Lock()
	h2o.oxygen++
	if h2o.hydrogen >= 2 {
		h2o.hydroQueue.Signal(2)
		h2o.hydrogen -= 2
		h2o.oxyQueue.Signal(1)
		h2o.oxygen--
	} else {
		h2o.mu.Unlock()
	}

	h2o.oxyQueue.Wait()
	bond()
	h2o.barrier.Wait()
	h2o.mu.Unlock()
}
