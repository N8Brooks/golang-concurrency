//go:build challenge

// Package babooncrossing contains the challenge version of the baboon crossing
// problem.
//
// Baboons cross a canyon on a single rope. Baboons going in opposite
// directions must never be on the rope at the same time, the rope may hold at
// most five baboons, and a continuing stream in one direction must not starve
// baboons waiting to go the other way.
package babooncrossing

import "context"

type Crossing struct{}

func NewCrossing() *Crossing {
	return &Crossing{}
}

func (c *Crossing) Left(ctx context.Context, arrive, cross, exit func()) error {
	panic("unimplemented")
}

func (c *Crossing) Right(ctx context.Context, arrive, cross, exit func()) error {
	panic("unimplemented")
}
