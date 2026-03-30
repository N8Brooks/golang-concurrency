package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/rendezvous/solutions"
	"github.com/N8Brooks/golang-concurrency/rendezvous/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.Rendezvous {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Rendezvous {
		return solutions.NewSemaphore()
	})
}
