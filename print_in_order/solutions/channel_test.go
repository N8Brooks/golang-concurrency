package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/print_in_order/solutions"
	"github.com/N8Brooks/golang-concurrency/print_in_order/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.PrintInOrder {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.PrintInOrder {
		return solutions.NewChannel()
	})
}
