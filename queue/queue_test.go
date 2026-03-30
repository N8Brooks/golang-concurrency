//go:build challenge

package queue_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/queue"
	"github.com/N8Brooks/golang-concurrency/queue/testsuite"
)

func TestQueue(t *testing.T) {
	testsuite.Run(t, func() testsuite.Queue {
		return queue.NewQueue()
	})
}

func BenchmarkQueue(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Queue {
		return queue.NewQueue()
	})
}
