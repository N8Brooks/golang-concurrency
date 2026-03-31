package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/throttle/solutions"
	"github.com/N8Brooks/golang-concurrency/throttle/testsuite"
)

func TestTimer(t *testing.T) {
	testsuite.Run(t, solutions.Throttle)
}
