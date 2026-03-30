package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/bounded_blocking_queue/solutions"
	"github.com/N8Brooks/golang-concurrency/bounded_blocking_queue/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.Queue {
		return solutions.NewChannel(capacity)
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.Queue {
		return solutions.NewChannel(capacity)
	})
}
