//go:build challenge

// Package nostarvereaderswriters contains the challenge version of the
// no-starve readers-writers problem.
//
// In the no-starve readers-writers problem, any number of readers may access
// the shared resource concurrently, but writers require exclusive access and
// newly arriving readers should not bypass a queued writer indefinitely.
package nostarvereaderswriters

import "context"

type NoStarveReadersWriters struct{}

func NewNoStarveReadersWriters() *NoStarveReadersWriters {
	return &NoStarveReadersWriters{}
}

func (rw *NoStarveReadersWriters) Reader(ctx context.Context, read func()) error {
	panic("unimplemented")
}

func (rw *NoStarveReadersWriters) Writer(ctx context.Context, write func()) error {
	panic("unimplemented")
}
