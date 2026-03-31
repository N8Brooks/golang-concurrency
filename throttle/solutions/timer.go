// Package solutions contains implementations of the Throttle problem.
package solutions

import (
	"sync"
	"time"
)

func Throttle(fn func(...int), d time.Duration) func(...int) {
	var mu sync.Mutex
	var pending []int
	var hasPending bool
	cooling := false

	var startTimer func()
	startTimer = func() {
		time.AfterFunc(d, func() {
			mu.Lock()
			if !hasPending {
				cooling = false
				mu.Unlock()
				return
			}
			args := append([]int(nil), pending...)
			hasPending = false
			mu.Unlock()

			fn(args...)
			startTimer()
		})
	}

	return func(args ...int) {
		mu.Lock()
		if !cooling {
			cooling = true
			mu.Unlock()

			fn(append([]int(nil), args...)...)
			startTimer()
			return
		}

		pending = append([]int(nil), args...)
		hasPending = true
		mu.Unlock()
	}
}
