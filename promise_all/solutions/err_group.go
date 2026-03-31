// Package solutions contains implementations of the Promise.all problem.
package solutions

import (
	"context"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	"golang.org/x/sync/errgroup"
)

func PromiseAllErrGroup[T any](ctx context.Context, functions []func(context.Context) promise.Promiser[T]) promise.Promiser[[]T] {
	return promise.New(func(resolve func([]T), reject func(error)) {
		g, ctx := errgroup.WithContext(ctx)
		values := make([]T, len(functions))

		for i, fn := range functions {
			g.Go(func() error {
				res, err := fn(ctx).Result()
				if err == nil {
					values[i] = res
				}
				return err
			})
		}

		if err := g.Wait(); err != nil {
			reject(err)
			return
		}

		resolve(values)
	})
}
