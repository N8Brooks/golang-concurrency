//go:build challenge

// Package barrier contains the challenge version of the barrier problem.
//
// In the barrier problem, a fixed number of participants each execute a first
// phase and a second phase. No participant may begin its second phase until
// every participant has completed its first phase.
//
// This challenge is the basic single-use barrier, not a reusable one.
package barrier

import "context"

type Barrier struct {
	parties int
}

func NewBarrier(parties int) *Barrier {
	return &Barrier{parties: parties}
}

func (b *Barrier) Wait(ctx context.Context, phase1, phase2 func()) error {
	panic("unimplemented")
}
