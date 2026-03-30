package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/river_crossing/solutions"
	"github.com/N8Brooks/golang-concurrency/river_crossing/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.RiverCrossing {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.RiverCrossing {
		return solutions.NewSemaphore()
	})
}
