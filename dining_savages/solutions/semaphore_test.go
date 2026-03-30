package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/dining_savages/solutions"
	"github.com/N8Brooks/golang-concurrency/dining_savages/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.DiningSavages {
		return solutions.NewSemaphore(capacity)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.DiningSavages {
		return solutions.NewSemaphore(capacity)
	})
}
