//go:build challenge

// Package hilzersbarbershop contains the challenge version of Hilzer's
// barbershop problem.
//
// The shop has three barbers, a four-seat sofa, room for twenty customers
// total, and a single cash register. Customers may have to stand before moving
// to the sofa, the sofa queue must be FIFO, up to three haircuts may happen at
// once, and payment must be serialized.
package hilzersbarbershop

import "context"

type Barbershop struct{}

func NewBarbershop() *Barbershop {
	return &Barbershop{}
}

func (s *Barbershop) Customer(ctx context.Context, enterShop, sitOnSofa, getHairCut, pay, exitShop, balk func()) error {
	panic("unimplemented")
}

func (s *Barbershop) Barber(ctx context.Context, cutHair, acceptPayment func()) error {
	panic("unimplemented")
}
