//go:build challenge

// Package readerswriters contains the challenge version of the readers-writers
// problem.
//
// In the readers-writers problem, any number of readers may access the shared
// resource concurrently, but writers require exclusive access.
package readerswriters

import "context"

type ReadersWriters struct{}

func NewReadersWriters() *ReadersWriters {
	return &ReadersWriters{}
}

func (rw *ReadersWriters) Reader(ctx context.Context, read func()) error {
	panic("unimplemented")
}

func (rw *ReadersWriters) Writer(ctx context.Context, write func()) error {
	panic("unimplemented")
}
