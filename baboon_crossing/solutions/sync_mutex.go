// Package solutions contains implementations of the baboon crossing problem.
package solutions

import (
	"context"
	"sync"
)

type direction int

const (
	none direction = iota
	left
	right
)

type waiter struct {
	ready chan struct{}
}

type SyncMutex struct {
	mu     sync.Mutex
	leftQ  []*waiter
	rightQ []*waiter
	active direction
	last   direction
	turn   direction
	onRope int
}

func NewSyncMutex() *SyncMutex {
	return &SyncMutex{}
}

func (c *SyncMutex) Left(ctx context.Context, arrive, cross, exit func()) error {
	return c.cross(ctx, left, arrive, cross, exit)
}

func (c *SyncMutex) Right(ctx context.Context, arrive, cross, exit func()) error {
	return c.cross(ctx, right, arrive, cross, exit)
}

func (c *SyncMutex) cross(ctx context.Context, dir direction, arrive, cross, exit func()) error {
	arrive()

	w := &waiter{ready: make(chan struct{})}

	c.mu.Lock()
	c.enqueueLocked(dir, w)
	if c.active != none && c.active != dir {
		c.turn = dir
	}
	c.maybeAdvanceLocked()
	c.mu.Unlock()

	select {
	case <-w.ready:
	case <-ctx.Done():
		c.mu.Lock()
		removed := c.removeLocked(dir, w)
		c.maybeAdvanceLocked()
		c.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-w.ready
	}

	cross()

	c.mu.Lock()
	c.onRope--
	if c.onRope == 0 {
		c.last = dir
		c.active = none
	}
	c.maybeAdvanceLocked()
	c.mu.Unlock()

	exit()
	return nil
}

func (c *SyncMutex) maybeAdvanceLocked() {
	if c.onRope > 0 {
		switch c.active {
		case left:
			if len(c.rightQ) > 0 {
				c.turn = right
				return
			}
			for c.onRope < 5 && len(c.leftQ) > 0 {
				w := c.popLocked(left)
				c.onRope++
				close(w.ready)
			}
		case right:
			if len(c.leftQ) > 0 {
				c.turn = left
				return
			}
			for c.onRope < 5 && len(c.rightQ) > 0 {
				w := c.popLocked(right)
				c.onRope++
				close(w.ready)
			}
		}
		return
	}

	next := c.chooseNextLocked()
	if next == none {
		return
	}

	c.active = next
	for c.onRope < 5 && c.queueLenLocked(next) > 0 {
		w := c.popLocked(next)
		c.onRope++
		close(w.ready)
	}

	if next == left && len(c.rightQ) > 0 {
		c.turn = right
	} else if next == right && len(c.leftQ) > 0 {
		c.turn = left
	} else {
		c.turn = none
	}
}

func (c *SyncMutex) chooseNextLocked() direction {
	if c.turn != none && c.queueLenLocked(c.turn) > 0 {
		return c.turn
	}

	switch {
	case len(c.leftQ) > 0 && len(c.rightQ) == 0:
		return left
	case len(c.rightQ) > 0 && len(c.leftQ) == 0:
		return right
	case len(c.leftQ) == 0 && len(c.rightQ) == 0:
		return none
	}

	if c.last == left {
		return right
	}
	if c.last == right {
		return left
	}
	return left
}

func (c *SyncMutex) enqueueLocked(dir direction, w *waiter) {
	if dir == left {
		c.leftQ = append(c.leftQ, w)
		return
	}
	c.rightQ = append(c.rightQ, w)
}

func (c *SyncMutex) popLocked(dir direction) *waiter {
	if dir == left {
		w := c.leftQ[0]
		c.leftQ = c.leftQ[1:]
		return w
	}
	w := c.rightQ[0]
	c.rightQ = c.rightQ[1:]
	return w
}

func (c *SyncMutex) queueLenLocked(dir direction) int {
	if dir == left {
		return len(c.leftQ)
	}
	return len(c.rightQ)
}

func (c *SyncMutex) removeLocked(dir direction, target *waiter) bool {
	var queue *[]*waiter
	if dir == left {
		queue = &c.leftQ
	} else {
		queue = &c.rightQ
	}

	for i, w := range *queue {
		if w != target {
			continue
		}
		copy((*queue)[i:], (*queue)[i+1:])
		(*queue)[len(*queue)-1] = nil
		*queue = (*queue)[:len(*queue)-1]
		return true
	}
	return false
}
