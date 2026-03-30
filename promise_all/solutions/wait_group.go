// Package solutions contains implementations of the Promise.all problem.
package solutions

import (
	"sync"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func PromiseAll[T any](functions []func() promise.Promiser[T]) promise.Promiser[[]T] {
	return promise.New(func(resolve func([]T), reject func(error)) {
		var wg sync.WaitGroup
		wg.Add(len(functions))

		values := make([]T, len(functions))

		for i, fn := range functions {
			go func(i int, fn func() promise.Promiser[T]) {
				defer wg.Done()
				if res, err := fn().Result(); err != nil {
					reject(err)
				} else {
					values[i] = res
				}
			}(i, fn)
		}

		wg.Wait()
		resolve(values)
	})
}
