//go:build challenge

package finitebufferproducerconsumer_test

import (
	"testing"

	finitebufferproducerconsumer "github.com/N8Brooks/golang-concurrency/finite_buffer_producer_consumer"
	"github.com/N8Brooks/golang-concurrency/finite_buffer_producer_consumer/testsuite"
)

func TestProducerConsumer(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.ProducerConsumer {
		return finitebufferproducerconsumer.NewProducerConsumer(capacity)
	})
}

func BenchmarkProducerConsumer(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.ProducerConsumer {
		return finitebufferproducerconsumer.NewProducerConsumer(capacity)
	})
}
