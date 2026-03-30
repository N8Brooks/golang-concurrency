//go:build challenge

// Package mutex contains the challenge version of the mutual exclusion problem.
//
// In the mutex problem, multiple threads need exclusive access to a critical
// section. At most one thread may hold the mutex at a time, and a waiting
// thread should be able to stop waiting if its context is canceled.
package mutex

import "context"

type Mutex struct{}

func NewMutex() *Mutex {
	return &Mutex{}
}

func (m *Mutex) Lock(ctx context.Context) error {
	panic("unimplemented")
}

func (m *Mutex) Unlock() {
	panic("unimplemented")
}
