package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/print_zero_even_odd/solutions"
	"github.com/N8Brooks/golang-concurrency/print_zero_even_odd/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.ZeroEvenOdd {
		return solutions.NewChannel(n)
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.ZeroEvenOdd {
		return solutions.NewChannel(n)
	})
}
