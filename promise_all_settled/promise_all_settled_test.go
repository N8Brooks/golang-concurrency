//go:build challenge

package promise_all_settled_test

import (
	"testing"

	promiseallsettled "github.com/N8Brooks/golang-concurrency/promise_all_settled"
	"github.com/N8Brooks/golang-concurrency/promise_all_settled/testsuite"
)

func TestPromiseAllSettled(t *testing.T) {
	testsuite.Run(t, promiseallsettled.PromiseAllSettled)
}

func BenchmarkPromiseAllSettled(b *testing.B) {
	testsuite.Benchmark(b, promiseallsettled.PromiseAllSettled)
}
