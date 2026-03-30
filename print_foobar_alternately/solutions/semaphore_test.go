package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/solutions"
	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.FooBar {
		return solutions.NewSemaphore(n)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.FooBar {
		return solutions.NewSemaphore(n)
	})
}
