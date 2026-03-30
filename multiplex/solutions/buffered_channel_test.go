package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/multiplex/solutions"
	"github.com/N8Brooks/golang-concurrency/multiplex/testsuite"
)

func TestBufferedChannel(t *testing.T) {
	testsuite.Run(t, func(limit int) testsuite.Multiplex {
		return solutions.NewBufferedChannel(limit)
	})
}

func BenchmarkBufferedChannel(b *testing.B) {
	testsuite.Benchmark(b, func(limit int) testsuite.Multiplex {
		return solutions.NewBufferedChannel(limit)
	})
}
