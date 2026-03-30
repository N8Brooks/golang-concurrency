//go:build challenge

// Package promise_time_limit contains the challenge version of the Promise Time
// Limit problem.
package promise_time_limit

import (
	"github.com/N8Brooks/golang-concurrency/internal/promise"
	"time"
)

func TimeLimit(fn func(...int) promise.Promiser[int], t time.Duration) func(...int) promise.Promiser[int] {
	panic("unimplemented")
}
