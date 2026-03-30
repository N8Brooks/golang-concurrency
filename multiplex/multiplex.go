//go:build challenge

// Package multiplex contains the challenge version of the multiplex problem.
//
// In the multiplex problem, up to a fixed number of callers may execute a
// critical section concurrently. Additional callers must wait until one of the
// active callers exits, unless their context is canceled first.
package multiplex

import "context"

type Multiplex struct{}

func NewMultiplex(limit int) *Multiplex {
	return &Multiplex{}
}

func (m *Multiplex) Run(ctx context.Context, criticalSection func()) {
	panic("unimplemented")
}
