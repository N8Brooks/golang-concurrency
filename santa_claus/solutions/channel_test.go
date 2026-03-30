package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/santa_claus/solutions"
	"github.com/N8Brooks/golang-concurrency/santa_claus/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.SantaClaus {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.SantaClaus {
		return solutions.NewChannel()
	})
}
