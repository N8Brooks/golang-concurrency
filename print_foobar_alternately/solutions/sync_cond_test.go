package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/solutions"
	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.FooBar {
		return solutions.NewSyncCond(n)
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.FooBar {
		return solutions.NewSyncCond(n)
	})
}
