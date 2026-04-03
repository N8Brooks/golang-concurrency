// Package solutions contains implementations of the river crossing problem.
package solutions

import (
	"context"
	"sync"
)

type semaphore struct {
	mu      sync.Mutex
	count   int
	waiters []chan struct{}
}

func (s *semaphore) Signal() {
	s.mu.Lock()
	if len(s.waiters) > 0 {
		waiter := s.waiters[0]
		copy(s.waiters, s.waiters[1:])
		s.waiters[len(s.waiters)-1] = nil
		s.waiters = s.waiters[:len(s.waiters)-1]
		s.mu.Unlock()
		close(waiter)
		return
	}
	s.count++
	s.mu.Unlock()
}

func (s *semaphore) Wait(ctx context.Context) bool {
	s.mu.Lock()
	if s.count > 0 {
		s.count--
		s.mu.Unlock()
		return true
	}
	if ctx.Err() != nil {
		s.mu.Unlock()
		return false
	}

	waiter := make(chan struct{})
	s.waiters = append(s.waiters, waiter)
	s.mu.Unlock()

	select {
	case <-waiter:
		return true
	case <-ctx.Done():
		s.mu.Lock()
		index := -1
		for i, candidate := range s.waiters {
			if candidate == waiter {
				index = i
				break
			}
		}
		if index >= 0 {
			copy(s.waiters[index:], s.waiters[index+1:])
			s.waiters[len(s.waiters)-1] = nil
			s.waiters = s.waiters[:len(s.waiters)-1]
			s.mu.Unlock()
			return false
		}
		s.mu.Unlock()
		<-waiter
		return true
	}
}

type Semaphore struct {
	mu      sync.Mutex
	hackers []*waiter
	serfs   []*waiter
}

func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

func (r *Semaphore) Hacker(ctx context.Context, board, rowBoat func()) error {
	return r.arrive(ctx, true, board, rowBoat)
}

func (r *Semaphore) Serf(ctx context.Context, board, rowBoat func()) error {
	return r.arrive(ctx, false, board, rowBoat)
}

func (r *Semaphore) arrive(ctx context.Context, hacker bool, board, rowBoat func()) error {
	w := &waiter{ready: make(chan assignment, 1)}

	r.mu.Lock()
	assignment, captain := r.enqueueAndSelect(hacker, w)
	if assignment == nil {
		r.mu.Unlock()

		select {
		case asg := <-w.ready:
			assignment = &asg
			captain = asg.captain
		case <-ctx.Done():
			r.mu.Lock()
			removed := r.removeWaiter(hacker, w)
			r.mu.Unlock()
			if removed {
				return ctx.Err()
			}
			asg := <-w.ready
			assignment = &asg
			captain = asg.captain
		}
	}

	board()
	assignment.boat.wait()

	if captain {
		rowBoat()
		r.mu.Unlock()
	}

	return nil
}

func (r *Semaphore) enqueueAndSelect(hacker bool, self *waiter) (*assignment, bool) {
	if hacker {
		r.hackers = append(r.hackers, self)
	} else {
		r.serfs = append(r.serfs, self)
	}

	switch {
	case len(r.hackers) >= 4:
		boat := newBoat()
		selected := r.hackers[:4]
		r.hackers = r.hackers[4:]
		return r.assignBoat(boat, self, selected, nil), true
	case len(r.serfs) >= 4:
		boat := newBoat()
		selected := r.serfs[:4]
		r.serfs = r.serfs[4:]
		return r.assignBoat(boat, self, nil, selected), true
	case len(r.hackers) >= 2 && len(r.serfs) >= 2:
		boat := newBoat()
		selectedHackers := r.hackers[:2]
		selectedSerfs := r.serfs[:2]
		r.hackers = r.hackers[2:]
		r.serfs = r.serfs[2:]
		return r.assignBoat(boat, self, selectedHackers, selectedSerfs), true
	default:
		return nil, false
	}
}

func (r *Semaphore) assignBoat(boat *boat, self *waiter, hackers, serfs []*waiter) *assignment {
	for _, waiter := range hackers {
		if waiter == self {
			continue
		}
		waiter.ready <- assignment{boat: boat}
	}
	for _, waiter := range serfs {
		if waiter == self {
			continue
		}
		waiter.ready <- assignment{boat: boat}
	}
	return &assignment{boat: boat, captain: true}
}

func (r *Semaphore) removeWaiter(hacker bool, target *waiter) bool {
	var queue []*waiter
	if hacker {
		queue = r.hackers
	} else {
		queue = r.serfs
	}

	for i, waiter := range queue {
		if waiter != target {
			continue
		}
		copy(queue[i:], queue[i+1:])
		queue[len(queue)-1] = nil
		queue = queue[:len(queue)-1]
		if hacker {
			r.hackers = queue
		} else {
			r.serfs = queue
		}
		return true
	}
	return false
}
