// Package solutions contains implementations of the Promise.allSettled problem.
package solutions

import (
	"sync"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	promiseallsettled "github.com/N8Brooks/golang-concurrency/promise_all_settled"
)

func PromiseAllSettled[T any](functions []func() promise.Promiser[T]) []promiseallsettled.Result[T] {
	var wg sync.WaitGroup
	wg.Add(len(functions))

	results := make([]promiseallsettled.Result[T], len(functions))

	for i, function := range functions {
		go func() {
			defer wg.Done()
			val, err := function().Result()
			results[i] = promiseallsettled.Result[T]{Val: val, Err: err}
		}()
	}

	wg.Wait()
	return results
}
