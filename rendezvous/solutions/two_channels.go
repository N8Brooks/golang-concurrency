// Package solutions contains implementations of the rendezvous problem.
package solutions

import "context"

type TwoChannels struct {
	aArrived chan struct{}
	bArrived chan struct{}
}

func NewTwoChannels() *TwoChannels {
	return &TwoChannels{
		aArrived: make(chan struct{}, 1),
		bArrived: make(chan struct{}, 1),
	}
}

func (r *TwoChannels) A(ctx context.Context, a1, a2 func()) {
	a1()
	select {
	case r.aArrived <- struct{}{}:
	case <-ctx.Done():
		return
	}
	select {
	case <-r.bArrived:
	case <-ctx.Done():
		return
	}
	a2()
}

func (r *TwoChannels) B(ctx context.Context, b1, b2 func()) {
	b1()
	select {
	case r.bArrived <- struct{}{}:
	case <-ctx.Done():
		return
	}
	select {
	case <-r.aArrived:
	case <-ctx.Done():
		return
	}
	b2()
}
