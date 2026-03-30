//go:build challenge

// Package diningsavages contains the challenge version of the dining savages
// problem.
//
// A fixed-size communal pot is shared by many savage threads and one cook
// thread. Savages may only take servings when the pot is non-empty, and the
// cook may only refill the pot once it is empty.
package diningsavages

import "context"

type DiningSavages struct {
	capacity int
}

func NewDiningSavages(capacity int) *DiningSavages {
	return &DiningSavages{capacity: capacity}
}

func (ds *DiningSavages) Savage(ctx context.Context, getServingFromPot, eat func()) {
	panic("unimplemented")
}

func (ds *DiningSavages) Cook(ctx context.Context, putServingsInPot func(servings int)) {
	panic("unimplemented")
}
