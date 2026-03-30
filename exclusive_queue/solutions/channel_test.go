package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/exclusive_queue/solutions"
	"github.com/N8Brooks/golang-concurrency/exclusive_queue/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExclusiveQueue {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExclusiveQueue {
		return solutions.NewChannel()
	})
}
