package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/fifo_barbershop/solutions"
	"github.com/N8Brooks/golang-concurrency/fifo_barbershop/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.FIFOBarbershop {
		return solutions.NewSemaphore(capacity)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.FIFOBarbershop {
		return solutions.NewSemaphore(capacity)
	})
}
