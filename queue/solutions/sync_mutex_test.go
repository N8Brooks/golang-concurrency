package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/queue/solutions"
	"github.com/N8Brooks/golang-concurrency/queue/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.Queue {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Queue {
		return solutions.NewSyncMutex()
	})
}
