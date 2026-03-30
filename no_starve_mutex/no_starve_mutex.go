//go:build challenge

// Package no_starve_mutex contains the challenge version of the bounded-waiting
// mutex problem.
//
// Multiple threads need mutual exclusion, but unlike a basic mutex, a thread
// that starts waiting must not be overtaken indefinitely by later arrivals.
// With a finite number of threads, the number of later threads that can
// proceed before it must be bounded.
package no_starve_mutex

import "context"

type NoStarveMutex struct{}

func NewNoStarveMutex() *NoStarveMutex {
	return &NoStarveMutex{}
}

func (m *NoStarveMutex) Lock(ctx context.Context) error {
	panic("unimplemented")
}

func (m *NoStarveMutex) Unlock() {
	panic("unimplemented")
}
