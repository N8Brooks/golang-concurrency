//go:build challenge

// Package unisexbathroom contains the challenge version of the unisex
// bathroom problem.
//
// In the unisex bathroom problem, men and women may use the bathroom, but
// never at the same time, and there may never be more than three people
// inside.
package unisexbathroom

import "context"

type UnisexBathroom struct{}

func NewUnisexBathroom() *UnisexBathroom {
	return &UnisexBathroom{}
}

func (b *UnisexBathroom) Male(ctx context.Context, bathroom func()) error {
	panic("unimplemented")
}

func (b *UnisexBathroom) Female(ctx context.Context, bathroom func()) error {
	panic("unimplemented")
}
