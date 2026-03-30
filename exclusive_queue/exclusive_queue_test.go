//go:build challenge

package exclusive_queue_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/exclusive_queue"
	"github.com/N8Brooks/golang-concurrency/exclusive_queue/testsuite"
)

func TestExclusiveQueue(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExclusiveQueue {
		return exclusive_queue.NewExclusiveQueue()
	})
}

func BenchmarkExclusiveQueue(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExclusiveQueue {
		return exclusive_queue.NewExclusiveQueue()
	})
}
