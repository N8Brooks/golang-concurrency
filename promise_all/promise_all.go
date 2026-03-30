//go:build challenge

// Package promise_all contains the challenge version of the Promise.all style
// problem.
package promise_all

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func PromiseAll[T any](functions []func() promise.Promiser[T]) promise.Promiser[[]T] {
	panic("unimplemented")
}
