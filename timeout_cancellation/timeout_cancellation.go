//go:build challenge

// Package timeout_cancellation contains the challenge version of the Timeout
// Cancellation problem.
package timeout_cancellation

import "time"

func Cancellable(fn func(...int), args []int, t time.Duration) func() {
	panic("unimplemented")
}
