// Package solutions contains implementations of the Timeout Cancellation problem.
package solutions

import (
	"sync"
	"time"
)

func Cancellable(fn func(...int), args []int, t time.Duration) func() {
	timer := time.AfterFunc(t, func() {
		fn(append([]int(nil), args...)...)
	})

	var once sync.Once
	return func() {
		once.Do(func() {
			timer.Stop()
		})
	}
}
