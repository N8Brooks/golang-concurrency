package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/no_starve_unisex_bathroom/solutions"
	"github.com/N8Brooks/golang-concurrency/no_starve_unisex_bathroom/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.NoStarveUnisexBathroom {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.NoStarveUnisexBathroom {
		return solutions.NewChannel()
	})
}
