//go:build challenge

// Package roller_coaster contains the challenge version of the roller coaster
// problem.
//
// Many passengers wait to take rides in a car with fixed capacity C. The car
// may only depart when exactly C passengers have boarded, and passengers may
// only unboard after the car has unloaded.
package roller_coaster

import "context"

type RollerCoaster struct{}

func NewRollerCoaster(capacity int) *RollerCoaster {
	return &RollerCoaster{}
}

func (r *RollerCoaster) Passenger(ctx context.Context, board, unboard func()) error {
	panic("unimplemented")
}

func (r *RollerCoaster) Car(ctx context.Context, load, run, unload func()) error {
	panic("unimplemented")
}
