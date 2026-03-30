//go:build challenge

// Package multi_car_roller_coaster contains the challenge version of the
// multi-car roller coaster problem.
//
// Multiple cars share a loading area and an unloading area. Only one car may
// board passengers at a time, but several cars may be running on the track
// concurrently. Because cars cannot pass one another, they must unload in the
// same cyclic order they boarded.
package multi_car_roller_coaster

import "context"

type MultiCarRollerCoaster struct{}

func NewMultiCarRollerCoaster(cars, capacity int) *MultiCarRollerCoaster {
	return &MultiCarRollerCoaster{}
}

func (r *MultiCarRollerCoaster) Passenger(ctx context.Context, board, unboard func()) error {
	panic("unimplemented")
}

func (r *MultiCarRollerCoaster) Car(ctx context.Context, car int, load, run, unload func()) error {
	panic("unimplemented")
}
