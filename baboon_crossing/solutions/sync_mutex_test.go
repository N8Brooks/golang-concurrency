package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/baboon_crossing/solutions"
	"github.com/N8Brooks/golang-concurrency/baboon_crossing/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.Crossing {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Crossing {
		return solutions.NewSyncMutex()
	})
}
