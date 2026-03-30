//go:build challenge

// Package fifo_barbershop contains the challenge version of the FIFO sleeping
// barber problem.
//
// A barbershop has room for a limited number of customers. If the shop is
// full, a newly arriving customer must balk instead of waiting. Otherwise the
// customer waits until the barber is ready, gets exactly one haircut, and then
// leaves. Unlike the basic barbershop problem, waiting customers must be
// served strictly in arrival order.
package fifo_barbershop

import "context"

type FIFOBarbershop struct{}

func NewFIFOBarbershop(capacity int) *FIFOBarbershop {
	return &FIFOBarbershop{}
}

func (b *FIFOBarbershop) Customer(ctx context.Context, getHairCut, balk func()) error {
	panic("unimplemented")
}

func (b *FIFOBarbershop) Barber(ctx context.Context, cutHair func()) error {
	panic("unimplemented")
}
