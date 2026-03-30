package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/building_h2o/solutions"
	"github.com/N8Brooks/golang-concurrency/building_h2o/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.H2O {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.H2O {
		return solutions.NewChannel()
	})
}
