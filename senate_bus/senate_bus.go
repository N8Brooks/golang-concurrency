//go:build challenge

// Package senatebus contains the challenge version of the Senate bus problem.
//
// In the Senate bus problem, riders wait for a bus with fixed capacity. When a
// bus arrives, it boards up to capacity from the riders already waiting,
// excludes late arrivals until the next bus, and departs when all selected
// riders have boarded.
package senatebus

import "context"

type SenateBus struct {
	capacity int
}

func NewSenateBus(capacity int) *SenateBus {
	return &SenateBus{capacity: capacity}
}

func (b *SenateBus) Rider(ctx context.Context, boardBus func()) error {
	panic("unimplemented")
}

func (b *SenateBus) Bus(ctx context.Context, depart func()) error {
	panic("unimplemented")
}
