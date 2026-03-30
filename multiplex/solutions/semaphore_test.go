package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/multiplex/solutions"
	"github.com/N8Brooks/golang-concurrency/multiplex/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(limit int) testsuite.Multiplex {
		return solutions.NewSemaphore(limit)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(limit int) testsuite.Multiplex {
		return solutions.NewSemaphore(limit)
	})
}
