package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/senate_bus/solutions"
	"github.com/N8Brooks/golang-concurrency/senate_bus/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.SenateBus {
		return solutions.NewChannel(capacity)
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.SenateBus {
		return solutions.NewChannel(capacity)
	})
}
