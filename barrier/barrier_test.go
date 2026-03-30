//go:build challenge

package barrier_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/barrier"
	"github.com/N8Brooks/golang-concurrency/barrier/testsuite"
)

func TestBarrier(t *testing.T) {
	testsuite.Run(t, func(parties int) testsuite.Barrier {
		return barrier.NewBarrier(parties)
	})
}

func BenchmarkBarrier(b *testing.B) {
	testsuite.Benchmark(b, func(parties int) testsuite.Barrier {
		return barrier.NewBarrier(parties)
	})
}
