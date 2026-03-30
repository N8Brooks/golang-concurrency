//go:build challenge

package dining_philosophers_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/dining_philosophers"
	"github.com/N8Brooks/golang-concurrency/dining_philosophers/testsuite"
)

func TestDiningPhilosophers(t *testing.T) {
	testsuite.Run(t, func() testsuite.DiningPhilosophers {
		return dining_philosophers.NewDiningPhilosophers()
	})
}

func BenchmarkDiningPhilosophers(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.DiningPhilosophers {
		return dining_philosophers.NewDiningPhilosophers()
	})
}
