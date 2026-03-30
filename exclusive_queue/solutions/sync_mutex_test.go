package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/exclusive_queue/solutions"
	"github.com/N8Brooks/golang-concurrency/exclusive_queue/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExclusiveQueue {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExclusiveQueue {
		return solutions.NewSyncMutex()
	})
}
