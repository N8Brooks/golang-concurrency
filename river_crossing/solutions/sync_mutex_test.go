package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/river_crossing/solutions"
	"github.com/N8Brooks/golang-concurrency/river_crossing/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.RiverCrossing {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.RiverCrossing {
		return solutions.NewSyncMutex()
	})
}
