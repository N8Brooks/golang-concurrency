package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/finite_buffer_producer_consumer/solutions"
	"github.com/N8Brooks/golang-concurrency/finite_buffer_producer_consumer/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.ProducerConsumer {
		return solutions.NewSyncCond(capacity)
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.ProducerConsumer {
		return solutions.NewSyncCond(capacity)
	})
}
