// Package solutions contains implementations of the callback-to-promise conversion problem.
package solutions

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func Promisify(fn func(callback func(int, error), args ...int)) func(args ...int) promise.Promiser[int] {
	return func(args ...int) promise.Promiser[int] {
		return promise.New(func(resolve func(int), reject func(error)) {
			fn(func(result int, err error) {
				if err != nil {
					reject(err)
					return
				}
				resolve(result)
			}, args...)
		})
	}
}
