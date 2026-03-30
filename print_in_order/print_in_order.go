//go:build challenge

// Package print_in_order contains the challenge version of the print-in-order
// problem.
//
// Three callers may invoke First, Second, and Third in any order, but the
// callbacks must still execute strictly in the order first, second, third.
package print_in_order

import "context"

type PrintInOrder struct{}

func NewPrintInOrder() *PrintInOrder {
	return &PrintInOrder{}
}

func (p *PrintInOrder) First(ctx context.Context, first func()) error {
	panic("unimplemented")
}

func (p *PrintInOrder) Second(ctx context.Context, second func()) error {
	panic("unimplemented")
}

func (p *PrintInOrder) Third(ctx context.Context, third func()) error {
	panic("unimplemented")
}
