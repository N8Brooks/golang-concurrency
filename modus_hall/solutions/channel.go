// Package solutions contains implementations of the Modus Hall problem.
package solutions

import (
	"context"
	"sync"
)

type side uint8

const (
	sideNeutral side = iota
	sideHeathen
	sidePrude
)

type waiter struct {
	ready chan struct{}
}

type Channel struct {
	mu              sync.Mutex
	ruling          side
	transition      side
	activeHeathens  int
	activePrudes    int
	waitingHeathens []*waiter
	waitingPrudes   []*waiter
}

func NewChannel() *Channel {
	return &Channel{}
}

func (mh *Channel) Heathen(ctx context.Context, cross func()) error {
	if err := mh.enter(ctx, sideHeathen); err != nil {
		return err
	}
	cross()
	mh.leave(sideHeathen)
	return nil
}

func (mh *Channel) Prude(ctx context.Context, cross func()) error {
	if err := mh.enter(ctx, sidePrude); err != nil {
		return err
	}
	cross()
	mh.leave(sidePrude)
	return nil
}

func (mh *Channel) enter(ctx context.Context, s side) error {
	mh.mu.Lock()
	if mh.canEnterLocked(s) {
		mh.admitImmediateLocked(s)
		mh.mu.Unlock()
		return nil
	}

	w := &waiter{ready: make(chan struct{})}
	if s == sideHeathen {
		mh.waitingHeathens = append(mh.waitingHeathens, w)
	} else {
		mh.waitingPrudes = append(mh.waitingPrudes, w)
	}
	mh.advanceLocked()
	mh.mu.Unlock()

	select {
	case <-w.ready:
		return nil
	case <-ctx.Done():
	}

	mh.mu.Lock()
	if mh.removeWaiterLocked(s, w) {
		mh.advanceLocked()
		mh.mu.Unlock()
		return ctx.Err()
	}
	mh.mu.Unlock()

	<-w.ready
	return nil
}

func (mh *Channel) leave(s side) {
	mh.mu.Lock()
	if s == sideHeathen {
		mh.activeHeathens--
	} else {
		mh.activePrudes--
	}
	mh.advanceLocked()
	mh.mu.Unlock()
}

func (mh *Channel) canEnterLocked(s side) bool {
	if mh.ruling == sideNeutral {
		return true
	}
	return mh.ruling == s && mh.transition == sideNeutral
}

func (mh *Channel) admitImmediateLocked(s side) {
	if mh.ruling == sideNeutral {
		mh.ruling = s
	}
	if s == sideHeathen {
		mh.activeHeathens++
	} else {
		mh.activePrudes++
	}
}

func (mh *Channel) advanceLocked() {
	switch mh.ruling {
	case sideNeutral:
		switch {
		case len(mh.waitingHeathens) > len(mh.waitingPrudes):
			mh.admitWaitingLocked(sideHeathen)
		case len(mh.waitingPrudes) > 0:
			mh.admitWaitingLocked(sidePrude)
		case len(mh.waitingHeathens) > 0:
			mh.admitWaitingLocked(sideHeathen)
		}
	case sideHeathen:
		if mh.activeHeathens == 0 {
			switch {
			case len(mh.waitingPrudes) > 0:
				mh.admitWaitingLocked(sidePrude)
			case len(mh.waitingHeathens) > 0:
				mh.admitWaitingLocked(sideHeathen)
			default:
				mh.ruling = sideNeutral
				mh.transition = sideNeutral
			}
			return
		}

		if len(mh.waitingPrudes) > mh.activeHeathens {
			mh.transition = sidePrude
		} else if mh.transition == sidePrude {
			mh.transition = sideNeutral
		}
	case sidePrude:
		if mh.activePrudes == 0 {
			switch {
			case len(mh.waitingHeathens) > 0:
				mh.admitWaitingLocked(sideHeathen)
			case len(mh.waitingPrudes) > 0:
				mh.admitWaitingLocked(sidePrude)
			default:
				mh.ruling = sideNeutral
				mh.transition = sideNeutral
			}
			return
		}

		if len(mh.waitingHeathens) > mh.activePrudes {
			mh.transition = sideHeathen
		} else if mh.transition == sideHeathen {
			mh.transition = sideNeutral
		}
	}
}

func (mh *Channel) admitWaitingLocked(s side) {
	var batch []*waiter

	if s == sideHeathen {
		batch = mh.waitingHeathens
		mh.waitingHeathens = nil
		mh.activeHeathens += len(batch)
		mh.activePrudes = 0
	} else {
		batch = mh.waitingPrudes
		mh.waitingPrudes = nil
		mh.activePrudes += len(batch)
		mh.activeHeathens = 0
	}

	mh.ruling = s
	mh.transition = sideNeutral

	for _, w := range batch {
		close(w.ready)
	}
}

func (mh *Channel) removeWaiterLocked(s side, target *waiter) bool {
	if s == sideHeathen {
		var removed bool
		mh.waitingHeathens, removed = removeWaiter(mh.waitingHeathens, target)
		return removed
	}

	var removed bool
	mh.waitingPrudes, removed = removeWaiter(mh.waitingPrudes, target)
	return removed
}

func removeWaiter(queue []*waiter, target *waiter) ([]*waiter, bool) {
	for i, w := range queue {
		if w == target {
			copy(queue[i:], queue[i+1:])
			return queue[:len(queue)-1], true
		}
	}
	return queue, false
}
