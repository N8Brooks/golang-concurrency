//go:build challenge

package promise_time_limit_test

import (
	"testing"

	promisetimelimit "github.com/N8Brooks/golang-concurrency/promise_time_limit"
	"github.com/N8Brooks/golang-concurrency/promise_time_limit/testsuite"
)

func TestTimeLimit(t *testing.T) {
	testsuite.Run(t, promisetimelimit.TimeLimit)
}

func BenchmarkTimeLimit(b *testing.B) {
	testsuite.Benchmark(b, promisetimelimit.TimeLimit)
}
