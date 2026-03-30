//go:build challenge

package boundedblockingqueue_test

import (
	"testing"

	boundedblockingqueue "github.com/N8Brooks/golang-concurrency/bounded_blocking_queue"
	"github.com/N8Brooks/golang-concurrency/bounded_blocking_queue/testsuite"
)

func TestBoundedBlockingQueue(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.Queue {
		return boundedblockingqueue.NewQueue(capacity)
	})
}

func BenchmarkBoundedBlockingQueue(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.Queue {
		return boundedblockingqueue.NewQueue(capacity)
	})
}
