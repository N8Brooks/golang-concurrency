package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/roller_coaster/solutions"
	"github.com/N8Brooks/golang-concurrency/roller_coaster/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.RollerCoaster {
		return solutions.NewSemaphore(capacity)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.RollerCoaster {
		return solutions.NewSemaphore(capacity)
	})
}
