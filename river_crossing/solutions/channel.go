package solutions

import (
	"context"
)

const (
	captain = iota
	passanger1
	passanger2
	passanger3
)

type Channel struct {
	hackerPair chan struct{}
	serfPair   chan struct{}
	pair1      chan struct{}
	pair2      chan struct{}
	x          chan struct{}
	y          chan struct{}
	z          chan struct{}
}

func NewChannel() *Channel {
	return &Channel{
		hackerPair: make(chan struct{}),
		serfPair:   make(chan struct{}),
		pair1:      make(chan struct{}),
		pair2:      make(chan struct{}),
		x:          make(chan struct{}),
		y:          make(chan struct{}),
		z:          make(chan struct{}),
	}
}

func (r *Channel) Hacker(ctx context.Context, board, rowBoat func()) error {
	return r.handshake(ctx, r.hackerPair, board, rowBoat)
}

func (r *Channel) Serf(ctx context.Context, board, rowBoat func()) error {
	return r.handshake(ctx, r.serfPair, board, rowBoat)
}

func (r *Channel) handshake(ctx context.Context, typeCh chan struct{}, board, rowBoat func()) error {
	var id uint8
	select {
	case <-ctx.Done():
		return ctx.Err()
	case typeCh <- struct{}{}:
	case <-typeCh:
		id |= 0b1
	}

	if id == captain {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case r.pair1 <- struct{}{}:
		case <-r.pair1:
			id |= 0b10
		}
	} else {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case r.pair2 <- struct{}{}:
		case <-r.pair2:
			id |= 0b10
		}
	}

	board()

	switch id {
	case captain:
		r.x <- struct{}{}
		r.z <- struct{}{}
		rowBoat()
	case passanger1:
		<-r.x
	case passanger2:
		r.y <- struct{}{}
		<-r.z
	case passanger3:
		<-r.y
	}

	return nil
}
