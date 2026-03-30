package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/barrier/solutions"
	"github.com/N8Brooks/golang-concurrency/barrier/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(parties int) testsuite.Barrier {
		return solutions.NewSemaphore(parties)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(parties int) testsuite.Barrier {
		return solutions.NewSemaphore(parties)
	})
}
