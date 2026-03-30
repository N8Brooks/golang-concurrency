package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/no_starve_readers_writers/solutions"
	"github.com/N8Brooks/golang-concurrency/no_starve_readers_writers/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.NoStarveReadersWriters {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.NoStarveReadersWriters {
		return solutions.NewChannel()
	})
}
