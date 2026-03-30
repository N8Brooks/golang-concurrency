package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/promise_all_settled/solutions"
	"github.com/N8Brooks/golang-concurrency/promise_all_settled/testsuite"
)

func TestWaitGroup(t *testing.T) {
	testsuite.Run(t, solutions.PromiseAllSettled)
}

func BenchmarkWaitGroup(b *testing.B) {
	testsuite.Benchmark(b, solutions.PromiseAllSettled)
}
