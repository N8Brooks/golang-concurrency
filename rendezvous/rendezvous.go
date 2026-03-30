//go:build challenge

// Package rendezvous contains the challenge version of the rendezvous problem.
//
// In the rendezvous problem, two threads each execute a first phase and a
// second phase. The synchronization requirement is that neither thread may
// begin its second phase until both threads have completed their first phase.
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
