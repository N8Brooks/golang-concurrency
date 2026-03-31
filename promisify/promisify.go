//go:build challenge

// Package promisify contains the challenge version of the callback-to-promise
// conversion problem.
package promisify

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func Promisify(fn func(callback func(int, error), args ...int)) func(args ...int) promise.Promiser[int] {
	panic("unimplemented")
}
