// Package solutions contains implementations of the Promise Pool problem.
package solutions

import (
	"sync"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func PromisePool(functions []func() promise.Promiser[any], n int) promise.Promiser[any] {
	ch := make(chan func() promise.Promiser[any], len(functions))

	var wg sync.WaitGroup
	wg.Add(n)

	for range n {
		go func() {
			defer wg.Done()
			for f := range ch {
				f().Result()
			}
		}()
	}

	for _, f := range functions {
		ch <- f
	}
	close(ch)

	wg.Wait()
	return promise.Resolve(any(struct{}{}))
}
