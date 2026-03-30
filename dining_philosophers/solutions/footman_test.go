package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/dining_philosophers/solutions"
	"github.com/N8Brooks/golang-concurrency/dining_philosophers/testsuite"
)

func TestFootman(t *testing.T) {
	testsuite.Run(t, func() testsuite.DiningPhilosophers {
		return solutions.NewFootman()
	})
}

func BenchmarkFootman(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.DiningPhilosophers {
		return solutions.NewFootman()
	})
}
