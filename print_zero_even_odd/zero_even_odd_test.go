//go:build challenge

package zeroevenodd_test

import (
	"testing"

	zeroevenodd "github.com/N8Brooks/golang-concurrency/print_zero_even_odd"
	"github.com/N8Brooks/golang-concurrency/print_zero_even_odd/testsuite"
)

func TestZeroEvenOdd(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.ZeroEvenOdd {
		return zeroevenodd.NewZeroEvenOdd(n)
	})
}

func BenchmarkZeroEvenOdd(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.ZeroEvenOdd {
		return zeroevenodd.NewZeroEvenOdd(n)
	})
}
