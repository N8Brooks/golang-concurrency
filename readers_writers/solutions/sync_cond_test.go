package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/readers_writers/solutions"
	"github.com/N8Brooks/golang-concurrency/readers_writers/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func() testsuite.ReadersWriters {
		return solutions.NewSyncCond()
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ReadersWriters {
		return solutions.NewSyncCond()
	})
}
