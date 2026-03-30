//go:build challenge

// Package faneuilhall contains the challenge version of the Faneuil Hall
// problem.
//
// Immigrants and spectators may enter the building only while the judge is
// absent. The judge may not confirm a ceremony until every immigrant who has
// entered has checked in. After confirmation, immigrants may collect their
// certificates, but they still may not leave until the judge departs.
package faneuilhall

import "context"

type FaneuilHall struct{}

func NewFaneuilHall() *FaneuilHall {
	return &FaneuilHall{}
}

func (fh *FaneuilHall) Immigrant(ctx context.Context, enter, checkIn, sitDown, swear, getCertificate, leave func()) error {
	panic("unimplemented")
}

func (fh *FaneuilHall) Judge(ctx context.Context, enter, confirm, leave func()) error {
	panic("unimplemented")
}

func (fh *FaneuilHall) Spectator(ctx context.Context, enter, spectate, leave func()) error {
	panic("unimplemented")
}
