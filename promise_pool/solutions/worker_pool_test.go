package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/promise_pool/solutions"
	"github.com/N8Brooks/golang-concurrency/promise_pool/testsuite"
)

func TestWorkerPool(t *testing.T) {
	testsuite.Run(t, solutions.PromisePool)
}

func BenchmarkWorkerPool(b *testing.B) {
	testsuite.Benchmark(b, solutions.PromisePool)
}
