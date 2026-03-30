//go:build challenge

// Package santaclaus contains the challenge version of the Santa Claus
// problem.
//
// Santa sleeps until either all nine reindeer have returned from vacation or a
// group of three elves needs help. Reindeer take priority over elves. After
// Santa helps one group of elves, the next group must wait until all three
// elves from the current group have finished getting help.
package santaclaus

import "context"

type SantaClaus struct{}

func NewSantaClaus() *SantaClaus {
	return &SantaClaus{}
}

func (sc *SantaClaus) Santa(ctx context.Context, prepareSleigh, helpElves func()) {
	panic("unimplemented")
}

func (sc *SantaClaus) Reindeer(ctx context.Context, getHitched func()) {
	panic("unimplemented")
}

func (sc *SantaClaus) Elf(ctx context.Context, getHelp func()) {
	panic("unimplemented")
}
