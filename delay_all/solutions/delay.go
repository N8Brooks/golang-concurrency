// Package solutions contains implementations of the Delay the Resolution of Each Promise problem.
package solutions

import (
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func DelayAll[T any](functions []func() promise.Promiser[T], d time.Duration) []func() promise.Promiser[T] {
	delayed := make([]func() promise.Promiser[T], len(functions))

	for i, fn := range functions {
		delayed[i] = func(fn func() promise.Promiser[T]) func() promise.Promiser[T] {
			return func() promise.Promiser[T] {
				return promise.New(func(resolve func(T), reject func(error)) {
					go func() {
						time.Sleep(d)
						val, err := fn().Result()
						if err != nil {
							reject(err)
							return
						}
						resolve(val)
					}()
				})
			}
		}(fn)
	}

	return delayed
}
