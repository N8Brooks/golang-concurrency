package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/extended_faneuil_hall/solutions"
	"github.com/N8Brooks/golang-concurrency/extended_faneuil_hall/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExtendedFaneuilHall {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExtendedFaneuilHall {
		return solutions.NewSemaphore()
	})
}
