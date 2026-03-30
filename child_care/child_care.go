//go:build challenge

// Package childcare contains the challenge version of the child care problem.
//
// In the child care problem, the center must maintain at least one adult for
// every three children who are currently inside. Adults may always enter, but
// children must wait for enough adult capacity and an adult who wants to leave
// may need to wait until the ratio is safe again.
package childcare

import "context"

type ChildCare struct{}

func NewChildCare() *ChildCare {
	return &ChildCare{}
}

func (c *ChildCare) Child(ctx context.Context, child func()) error {
	panic("unimplemented")
}

func (c *ChildCare) Adult(ctx context.Context, adult func()) error {
	panic("unimplemented")
}
