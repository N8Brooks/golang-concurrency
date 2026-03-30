//go:build challenge

// Package modushall contains the challenge version of the Modus Hall problem.
//
// Heathens and prudes may each cross the field concurrently with members of
// their own faction, but the two factions may not cross at the same time.
// Control of the field follows majority rule: once the waiting opposition
// outnumbers the currently active faction, new entrants from the active side
// are barred and control flips when the field clears.
package modushall

import "context"

type ModusHall struct{}

func NewModusHall() *ModusHall {
	return &ModusHall{}
}

func (mh *ModusHall) Heathen(ctx context.Context, cross func()) error {
	panic("unimplemented")
}

func (mh *ModusHall) Prude(ctx context.Context, cross func()) error {
	panic("unimplemented")
}
