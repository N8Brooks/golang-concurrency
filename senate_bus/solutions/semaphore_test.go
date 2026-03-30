package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/senate_bus/solutions"
	"github.com/N8Brooks/golang-concurrency/senate_bus/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.SenateBus {
		return solutions.NewSemaphore(capacity)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.SenateBus {
		return solutions.NewSemaphore(capacity)
	})
}
