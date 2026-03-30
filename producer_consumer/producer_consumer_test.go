//go:build challenge

package producerconsumer_test

import (
	"testing"

	producerconsumer "github.com/N8Brooks/golang-concurrency/producer_consumer"
	"github.com/N8Brooks/golang-concurrency/producer_consumer/testsuite"
)

func TestProducerConsumer(t *testing.T) {
	testsuite.Run(t, func() testsuite.ProducerConsumer {
		return producerconsumer.NewProducerConsumer()
	})
}

func BenchmarkProducerConsumer(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ProducerConsumer {
		return producerconsumer.NewProducerConsumer()
	})
}
