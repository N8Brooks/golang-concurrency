package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/custom_interval/solutions"
	"github.com/N8Brooks/golang-concurrency/custom_interval/testsuite"
)

func TestIntervals(t *testing.T) {
	testsuite.Run(t, solutions.CustomInterval, solutions.CustomClearInterval)
}
