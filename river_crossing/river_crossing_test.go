//go:build challenge

package rivercrossing_test

import (
	"testing"

	rivercrossing "github.com/N8Brooks/golang-concurrency/river_crossing"
	"github.com/N8Brooks/golang-concurrency/river_crossing/testsuite"
)

func TestRiverCrossing(t *testing.T) {
	testsuite.Run(t, func() testsuite.RiverCrossing {
		return rivercrossing.NewRiverCrossing()
	})
}

func BenchmarkRiverCrossing(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.RiverCrossing {
		return rivercrossing.NewRiverCrossing()
	})
}
