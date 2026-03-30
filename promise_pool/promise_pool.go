//go:build challenge

// Package promise_pool contains the challenge version of the Promise Pool
// problem.
package promise_pool

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func PromisePool(functions []func() promise.Promiser[any], n int) promise.Promiser[any] {
	panic("unimplemented")
}
