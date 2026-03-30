package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/unisex_bathroom/solutions"
	"github.com/N8Brooks/golang-concurrency/unisex_bathroom/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.UnisexBathroom {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.UnisexBathroom {
		return solutions.NewChannel()
	})
}
