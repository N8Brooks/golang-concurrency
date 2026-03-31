// Package solutions contains implementations of the Custom Interval problem.
package solutions

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	nextID atomic.Int64

	registryMu sync.Mutex
	registry   = make(map[int]chan struct{})
)

func CustomInterval(fn func(), delay, period time.Duration) int {
	id := int(nextID.Add(1))
	done := make(chan struct{})

	registryMu.Lock()
	registry[id] = done
	registryMu.Unlock()

	go func() {
		next := delay
		count := 0

		for {
			select {
			case <-done:
				return
			case <-time.After(next):
			}

			select {
			case <-done:
				return
			default:
			}

			fn()
			count++
			next = delay + period*time.Duration(count)
		}
	}()

	return id
}

func CustomClearInterval(id int) {
	registryMu.Lock()
	done, ok := registry[id]
	if ok {
		delete(registry, id)
	}
	registryMu.Unlock()

	if ok {
		close(done)
	}
}
