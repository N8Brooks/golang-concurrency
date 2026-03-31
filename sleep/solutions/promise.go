// Package solutions contains implementations of the Sleep problem.
package solutions

import (
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func Sleep(d time.Duration) promise.Promiser[struct{}] {
	return promise.New(func(resolve func(struct{}), reject func(error)) {
		go func() {
			time.Sleep(d)
			resolve(struct{}{})
		}()
	})
}
