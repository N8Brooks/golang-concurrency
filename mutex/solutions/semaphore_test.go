package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/mutex/solutions"
	"github.com/N8Brooks/golang-concurrency/mutex/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.Mutex {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Mutex {
		return solutions.NewSemaphore()
	})
}
