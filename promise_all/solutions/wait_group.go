// Package solutions contains implementations of the Promise.all problem.
package solutions

import (
	"context"
	"sync"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func PromiseAllWaitGroup[T any](ctx context.Context, functions []func(context.Context) promise.Promiser[T]) promise.Promiser[[]T] {
	return promise.New(func(resolve func([]T), reject func(error)) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(len(functions))

		values := make([]T, len(functions))

		for i, fn := range functions {
			go func() {
				defer wg.Done()
				if res, err := fn(ctx).ResultContext(ctx); err != nil {
					cancel()
					reject(err)
				} else {
					values[i] = res
				}
			}()
		}

		wg.Wait()
		resolve(values)
	})
}
