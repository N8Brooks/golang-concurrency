// Package solutions contains implementations of the rendezvous problem.
package solutions

import "context"

type OneChannel struct {
	handshake chan struct{}
}

func NewOneChannel() *OneChannel {
	return &OneChannel{
		handshake: make(chan struct{}),
	}
}

func (r *OneChannel) A(ctx context.Context, a1, a2 func()) {
	a1()
	select {
	case <-ctx.Done():
		return
	case r.handshake <- struct{}{}:
	}
	a2()
}

func (r *OneChannel) B(ctx context.Context, b1, b2 func()) {
	b1()
	select {
	case <-ctx.Done():
		return
	case <-r.handshake:
	}
	b2()
}
