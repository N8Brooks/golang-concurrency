package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/interval_cancellation/solutions"
	"github.com/N8Brooks/golang-concurrency/interval_cancellation/testsuite"
)

func TestTicker(t *testing.T) {
	testsuite.Run(t, solutions.Cancellable)
}
