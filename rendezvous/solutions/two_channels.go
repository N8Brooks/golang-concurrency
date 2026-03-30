package solutions

import "context"

type TwoChannel struct {
	aArrived chan struct{}
	bArrived chan struct{}
}

func NewTwoChannel() *TwoChannel {
	return &TwoChannel{
		aArrived: make(chan struct{}, 1),
		bArrived: make(chan struct{}, 1),
	}
}

func (r *TwoChannel) A(ctx context.Context, a1, a2 func()) {
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

func (r *TwoChannel) B(ctx context.Context, b1, b2 func()) {
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
