package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/promise_time_limit/solutions"
	"github.com/N8Brooks/golang-concurrency/promise_time_limit/testsuite"
)

func TestSelect(t *testing.T) {
	testsuite.Run(t, solutions.TimeLimit)
}

func BenchmarkSelect(b *testing.B) {
	testsuite.Benchmark(b, solutions.TimeLimit)
}
