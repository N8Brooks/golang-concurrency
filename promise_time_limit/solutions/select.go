// Package solutions contains implementations of the Promise Time Limit problem.
package solutions

import (
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	promisetimelimit "github.com/N8Brooks/golang-concurrency/promise_time_limit"
)

func TimeLimit(fn func(...int) promise.Promiser[int], t time.Duration) func(...int) promise.Promiser[int] {
	return func(args ...int) promise.Promiser[int] {
		resultCh := make(chan struct {
			val int
			err error
		}, 1)

		go func() {
			val, err := fn(args...).Result()
			resultCh <- struct {
				val int
				err error
			}{val: val, err: err}
		}()

		select {
		case <-time.After(t):
			return promise.Reject[int](promisetimelimit.ErrTimeLimitExceeded)
		case result := <-resultCh:
			if result.err != nil {
				return promise.Reject[int](result.err)
			}
			return promise.Resolve(result.val)
		}
	}
}
