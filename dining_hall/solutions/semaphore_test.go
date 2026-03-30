package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/dining_hall/solutions"
	"github.com/N8Brooks/golang-concurrency/dining_hall/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.DiningHall {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.DiningHall {
		return solutions.NewSemaphore()
	})
}
