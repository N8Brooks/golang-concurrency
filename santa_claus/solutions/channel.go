// Package solutions contains implementations of the Santa Claus problem.
package solutions

import (
	"context"
	"sync"
)

const numReindeer = 9

type requestState uint8

const (
	stateWaiting requestState = iota
	stateAdmitted
	stateCanceled
)

type reindeerRequest struct {
	release chan struct{}
	state   requestState
}

type elfRequest struct {
	release chan struct{}
	state   requestState
}

type Channel struct {
	mu          sync.Mutex
	wakeSanta   chan struct{}
	reindeer    []*reindeerRequest
	elves       []*elfRequest
	activeElves int
}

func NewChannel() *Channel {
	return &Channel{
		wakeSanta: make(chan struct{}, 1),
	}
}

func (sc *Channel) Santa(ctx context.Context, prepareSleigh, helpElves func()) {
	for {
		reindeerBatch, elfBatch := sc.nextBatch()
		switch {
		case len(reindeerBatch) != 0:
			prepareSleigh()
			sc.releaseReindeer(reindeerBatch)
		case len(elfBatch) != 0:
			helpElves()
		default:
			select {
			case <-ctx.Done():
				return
			case <-sc.wakeSanta:
			}
		}
	}
}

func (sc *Channel) Reindeer(ctx context.Context, getHitched func()) {
	req := &reindeerRequest{release: make(chan struct{})}

	sc.mu.Lock()
	sc.reindeer = append(sc.reindeer, req)
	ready := len(sc.reindeer) >= numReindeer
	sc.mu.Unlock()

	if ready {
		sc.notifySanta()
	}

	if !sc.waitForReindeer(ctx, req) {
		return
	}

	getHitched()
}

func (sc *Channel) Elf(ctx context.Context, getHelp func()) {
	req := &elfRequest{release: make(chan struct{})}

	sc.mu.Lock()
	sc.elves = append(sc.elves, req)
	ready := sc.activeElves == 0 && len(sc.elves) >= 3
	sc.mu.Unlock()

	if ready {
		sc.notifySanta()
	}

	if !sc.waitForElf(ctx, req) {
		return
	}

	getHelp()

	sc.mu.Lock()
	sc.activeElves--
	ready = sc.activeElves == 0 && len(sc.elves) >= 3
	sc.mu.Unlock()

	if ready {
		sc.notifySanta()
	}
}

func (sc *Channel) nextBatch() ([]*reindeerRequest, []*elfRequest) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if len(sc.reindeer) >= numReindeer {
		batch := append([]*reindeerRequest(nil), sc.reindeer[:numReindeer]...)
		sc.reindeer = sc.reindeer[numReindeer:]
		for _, req := range batch {
			req.state = stateAdmitted
		}
		return batch, nil
	}

	if sc.activeElves == 0 && len(sc.elves) >= 3 {
		batch := append([]*elfRequest(nil), sc.elves[:3]...)
		sc.elves = sc.elves[3:]
		sc.activeElves = 3
		for _, req := range batch {
			req.state = stateAdmitted
			close(req.release)
		}
		return nil, batch
	}

	return nil, nil
}

func (sc *Channel) releaseReindeer(batch []*reindeerRequest) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for _, req := range batch {
		close(req.release)
	}
}

func (sc *Channel) waitForReindeer(ctx context.Context, req *reindeerRequest) bool {
	select {
	case <-req.release:
		return true
	case <-ctx.Done():
	}

	sc.mu.Lock()
	if req.state == stateWaiting {
		sc.reindeer = removeReindeer(sc.reindeer, req)
		req.state = stateCanceled
		sc.mu.Unlock()
		return false
	}

	sc.mu.Unlock()
	<-req.release
	return true
}

func (sc *Channel) waitForElf(ctx context.Context, req *elfRequest) bool {
	select {
	case <-req.release:
		return true
	case <-ctx.Done():
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	if req.state == stateWaiting {
		sc.elves = removeElves(sc.elves, req)
		req.state = stateCanceled
		return false
	}

	return true
}

func (sc *Channel) notifySanta() {
	select {
	case sc.wakeSanta <- struct{}{}:
	default:
	}
}

func removeReindeer(queue []*reindeerRequest, target *reindeerRequest) []*reindeerRequest {
	for i, req := range queue {
		if req == target {
			copy(queue[i:], queue[i+1:])
			return queue[:len(queue)-1]
		}
	}
	return queue
}

func removeElves(queue []*elfRequest, target *elfRequest) []*elfRequest {
	for i, req := range queue {
		if req == target {
			copy(queue[i:], queue[i+1:])
			return queue[:len(queue)-1]
		}
	}
	return queue
}
