//go:build challenge

package promise_all_test

import (
	"testing"

	promiseall "github.com/N8Brooks/golang-concurrency/promise_all"
	"github.com/N8Brooks/golang-concurrency/promise_all/testsuite"
)

func TestPromiseAll(t *testing.T) {
	testsuite.Run(t, promiseall.PromiseAll[int])
}

func BenchmarkPromiseAll(b *testing.B) {
	testsuite.Benchmark(b, promiseall.PromiseAll[int])
}
