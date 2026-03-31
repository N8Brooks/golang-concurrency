//go:build challenge

// Package delay_all contains the challenge version of the Delay the Resolution
// of Each Promise problem.
package delay_all

import (
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func DelayAll[T any](functions []func() promise.Promiser[T], d time.Duration) []func() promise.Promiser[T] {
	panic("unimplemented")
}
