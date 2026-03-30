package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/extended_dining_hall/solutions"
	"github.com/N8Brooks/golang-concurrency/extended_dining_hall/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExtendedDiningHall {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExtendedDiningHall {
		return solutions.NewSyncMutex()
	})
}
