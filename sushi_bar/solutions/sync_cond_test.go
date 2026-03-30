package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/sushi_bar/solutions"
	"github.com/N8Brooks/golang-concurrency/sushi_bar/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.SushiBar {
		return solutions.NewSyncCond(capacity)
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.SushiBar {
		return solutions.NewSyncCond(capacity)
	})
}
