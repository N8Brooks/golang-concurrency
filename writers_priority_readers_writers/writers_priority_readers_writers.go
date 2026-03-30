//go:build challenge

// Package writerspriorityreaderswriters contains the challenge version of the
// writers-priority readers-writers problem.
//
// In the writers-priority readers-writers problem, any number of readers may
// access the shared resource concurrently, but writers require exclusive
// access and no new readers should enter while writers are queued.
package writerspriorityreaderswriters

import "context"

type WritersPriorityReadersWriters struct{}

func NewWritersPriorityReadersWriters() *WritersPriorityReadersWriters {
	return &WritersPriorityReadersWriters{}
}

func (rw *WritersPriorityReadersWriters) Reader(ctx context.Context, read func()) error {
	panic("unimplemented")
}

func (rw *WritersPriorityReadersWriters) Writer(ctx context.Context, write func()) error {
	panic("unimplemented")
}
