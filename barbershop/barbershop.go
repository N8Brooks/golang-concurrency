//go:build challenge

// Package barbershop contains the challenge version of the sleeping barber
// problem.
//
// A barbershop has room for a limited number of customers. If the shop is
// full, a newly arriving customer must balk instead of waiting. Otherwise the
// customer waits until the barber is ready, gets exactly one haircut, and then
// leaves. The barber sleeps when no customers are present.
package barbershop

import "context"

type Barbershop struct{}

func NewBarbershop(capacity int) *Barbershop {
	return &Barbershop{}
}

func (b *Barbershop) Customer(ctx context.Context, getHairCut, balk func()) error {
	panic("unimplemented")
}

func (b *Barbershop) Barber(ctx context.Context, cutHair func()) error {
	panic("unimplemented")
}
