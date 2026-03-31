// Package solutions contains implementations of the Interval Cancellation problem.
package solutions

import (
	"sync"
	"time"
)

func Cancellable(fn func(...int), args []int, t time.Duration) func() {
	done := make(chan struct{})
	var once sync.Once

	fn(append([]int(nil), args...)...)

	go func() {
		ticker := time.NewTicker(t)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fn(append([]int(nil), args...)...)
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(done)
		})
	}
}
