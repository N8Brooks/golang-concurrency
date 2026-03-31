package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/promise_all/solutions"
	"github.com/N8Brooks/golang-concurrency/promise_all/testsuite"
)

func TestWaitGroup(t *testing.T) {
	testsuite.Run(t, solutions.PromiseAllWaitGroup[int])
}

func BenchmarkWaitGroup(b *testing.B) {
	testsuite.Benchmark(b, solutions.PromiseAllWaitGroup[int])
}
