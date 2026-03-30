//go:build challenge

// Package extendedchildcare contains the challenge version of the extended
// child care problem.
//
// In the extended child care problem, the center must maintain at least one
// adult for every three children, but an adult that is waiting to leave still
// counts as present until it actually exits. That means children should still
// be able to enter whenever doing so is legal under the current number of
// adults inside.
package extendedchildcare

import "context"

type ExtendedChildCare struct{}

func NewExtendedChildCare() *ExtendedChildCare {
	return &ExtendedChildCare{}
}

func (c *ExtendedChildCare) Child(ctx context.Context, child func()) error {
	panic("unimplemented")
}

func (c *ExtendedChildCare) Adult(ctx context.Context, adult func()) error {
	panic("unimplemented")
}
