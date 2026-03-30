//go:build challenge

// Package extendedfaneuilhall contains the challenge version of the extended
// Faneuil Hall problem.
//
// Immigrants and spectators may enter the building only while the judge is
// absent. The judge may not confirm a ceremony until every immigrant who has
// entered has checked in. After confirmation, immigrants may collect their
// certificates, but they still may not leave until the judge departs. In the
// extended problem, the next judge must also wait until all sworn immigrants
// from the previous ceremony have left the building.
package extendedfaneuilhall

import "context"

type ExtendedFaneuilHall struct{}

func NewExtendedFaneuilHall() *ExtendedFaneuilHall {
	return &ExtendedFaneuilHall{}
}

func (fh *ExtendedFaneuilHall) Immigrant(ctx context.Context, enter, checkIn, sitDown, swear, getCertificate, leave func()) error {
	panic("unimplemented")
}

func (fh *ExtendedFaneuilHall) Judge(ctx context.Context, enter, confirm, leave func()) error {
	panic("unimplemented")
}

func (fh *ExtendedFaneuilHall) Spectator(ctx context.Context, enter, spectate, leave func()) error {
	panic("unimplemented")
}
