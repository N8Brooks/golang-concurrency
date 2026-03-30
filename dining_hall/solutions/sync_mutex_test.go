package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/dining_hall/solutions"
	"github.com/N8Brooks/golang-concurrency/dining_hall/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.DiningHall {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.DiningHall {
		return solutions.NewSyncMutex()
	})
}
