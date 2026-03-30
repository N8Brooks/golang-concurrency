//go:build challenge

// Package nostarveunisexbathroom contains the challenge version of the
// no-starve unisex bathroom problem.
//
// In the no-starve unisex bathroom problem, men and women may use the
// bathroom, but never at the same time, there may never be more than three
// people inside, and one sex should not be able to starve the other.
package nostarveunisexbathroom

import "context"

type NoStarveUnisexBathroom struct{}

func NewNoStarveUnisexBathroom() *NoStarveUnisexBathroom {
	return &NoStarveUnisexBathroom{}
}

func (b *NoStarveUnisexBathroom) Male(ctx context.Context, bathroom func()) error {
	panic("unimplemented")
}

func (b *NoStarveUnisexBathroom) Female(ctx context.Context, bathroom func()) error {
	panic("unimplemented")
}
