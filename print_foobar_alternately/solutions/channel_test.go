package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/solutions"
	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.FooBar {
		return solutions.NewChannel(n)
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.FooBar {
		return solutions.NewChannel(n)
	})
}
