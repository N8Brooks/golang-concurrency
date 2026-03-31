package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/timeout_cancellation/solutions"
	"github.com/N8Brooks/golang-concurrency/timeout_cancellation/testsuite"
)

func TestTimer(t *testing.T) {
	testsuite.Run(t, solutions.Cancellable)
}
