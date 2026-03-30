package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/producer_consumer/solutions"
	"github.com/N8Brooks/golang-concurrency/producer_consumer/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func() testsuite.ProducerConsumer {
		return solutions.NewSyncCond()
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ProducerConsumer {
		return solutions.NewSyncCond()
	})
}
