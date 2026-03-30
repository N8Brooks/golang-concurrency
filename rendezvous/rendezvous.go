//go:build challenge

package rendezvous

import "context"

type Rendezvous struct{}

func NewRendezvous() *Rendezvous {
	return &Rendezvous{}
}

func (r *Rendezvous) A(ctx context.Context, a1, a2 func()) {
	panic("unimplemented")
}

func (r *Rendezvous) B(ctx context.Context, b1, b2 func()) {
	panic("unimplemented")
}
