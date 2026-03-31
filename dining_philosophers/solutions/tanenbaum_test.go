package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/dining_philosophers/solutions"
	"github.com/N8Brooks/golang-concurrency/dining_philosophers/testsuite"
)

// NOTE: This solution fails the test suite because it starves
// func TestTanenbaum(t *testing.T) {
// 	testsuite.Run(t, func() testsuite.DiningPhilosophers {
// 		return solutions.NewTanenbaum()
// 	})
// }

func BenchmarkTanenbaum(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.DiningPhilosophers {
		return solutions.NewTanenbaum()
	})
}
