//go:build challenge

package promise_pool_test

import (
	"testing"

	promisepool "github.com/N8Brooks/golang-concurrency/promise_pool"
	"github.com/N8Brooks/golang-concurrency/promise_pool/testsuite"
)

func TestPromisePool(t *testing.T) {
	testsuite.Run(t, promisepool.PromisePool)
}

func BenchmarkPromisePool(b *testing.B) {
	testsuite.Benchmark(b, promisepool.PromisePool)
}
