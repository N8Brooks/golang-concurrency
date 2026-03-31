//go:build challenge

// Package sleep contains the challenge version of the Sleep problem.
package sleep

import (
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

func Sleep(d time.Duration) promise.Promiser[struct{}] {
	panic("unimplemented")
}
