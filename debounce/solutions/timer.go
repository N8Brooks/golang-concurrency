// Package solutions contains implementations of the debounce problem.
package solutions

import (
	"sync/atomic"
	"time"
)

func Debounce(fn func(...int), t time.Duration) func(...int) {
	var timer atomic.Pointer[time.Timer]
	var generation atomic.Uint64

	return func(args ...int) {
		next := generation.Add(1)
		captured := append([]int(nil), args...)

		callback := func() {
			if generation.Load() == next {
				fn(captured...)
			}
		}

		newTimer := time.AfterFunc(t, callback)
		oldTimer := timer.Swap(newTimer)
		if oldTimer != nil {
			oldTimer.Stop()
		}
	}
}
