package solutions_test

import (
	"testing"

	rivercrossing "github.com/N8Brooks/golang-concurrency/river_crossing/solutions"
	"github.com/N8Brooks/golang-concurrency/river_crossing/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.RiverCrossing {
		return rivercrossing.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.RiverCrossing {
		return rivercrossing.NewChannel()
	})
}
