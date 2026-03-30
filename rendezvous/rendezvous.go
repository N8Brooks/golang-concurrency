package rendezvous

import "context"

type Rendezvous struct {
	handshake chan struct{}
}

func NewRendezvous() *Rendezvous {
	return &Rendezvous{
		handshake: make(chan struct{}),
	}
}

func (r *Rendezvous) A(ctx context.Context, a1, a2 func()) {
	a1()
	select {
	case <-ctx.Done():
		return
	case r.handshake <- struct{}{}:
	}
	a2()
}

func (r *Rendezvous) B(ctx context.Context, b1, b2 func()) {
	b1()
	select {
	case <-ctx.Done():
		return
	case <-r.handshake:
	}
	b2()
}
