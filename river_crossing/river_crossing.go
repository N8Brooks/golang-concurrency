//go:build challenge

// Package rivercrossing contains the challenge version of the river crossing
// problem.
//
// Hackers and serfs arrive independently. A boat may leave only with four
// hackers, four serfs, or two hackers and two serfs. All four selected
// passengers must board before any subsequent boatload boards, and exactly one
// passenger in each boatload must row.
package rivercrossing

import (
	"context"
	"sync"
)

type RiverCrossing struct {
	mu      sync.Mutex
	hackers []*waiter
	serfs   []*waiter
}

func NewRiverCrossing() *RiverCrossing {
	return &RiverCrossing{}
}

func (r *RiverCrossing) Hacker(ctx context.Context, board, rowBoat func()) error {
	return r.arrive(ctx, true, board, rowBoat)
}

func (r *RiverCrossing) Serf(ctx context.Context, board, rowBoat func()) error {
	return r.arrive(ctx, false, board, rowBoat)
}

func (r *RiverCrossing) arrive(ctx context.Context, hacker bool, board, rowBoat func()) error {
	w := &waiter{
		ready: make(chan assignment, 1),
	}

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

func (r *RiverCrossing) enqueueAndSelect(hacker bool, self *waiter) (*assignment, bool) {
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

func (r *RiverCrossing) assignBoat(boat *boat, self *waiter, hackers, serfs []*waiter) *assignment {
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
	return &assignment{
		boat:    boat,
		captain: true,
	}
}

func (r *RiverCrossing) removeWaiter(hacker bool, target *waiter) bool {
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

type waiter struct {
	ready chan assignment
}

type assignment struct {
	boat    *boat
	captain bool
}

type boat struct {
	barrier *barrier
}

func newBoat() *boat {
	return &boat{
		barrier: newBarrier(4),
	}
}

func (b *boat) wait() {
	b.barrier.wait()
}

type barrier struct {
	total     int
	count     int
	mu        sync.Mutex
	turnstile chan struct{}
}

func newBarrier(total int) *barrier {
	return &barrier{
		total:     total,
		turnstile: make(chan struct{}, 1),
	}
}

func (b *barrier) wait() {
	b.mu.Lock()
	b.count++
	last := b.count == b.total
	b.mu.Unlock()

	if last {
		b.turnstile <- struct{}{}
	}

	<-b.turnstile
	b.turnstile <- struct{}{}
}
