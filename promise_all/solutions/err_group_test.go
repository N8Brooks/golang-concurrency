package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/promise_all/solutions"
	"github.com/N8Brooks/golang-concurrency/promise_all/testsuite"
)

func TestErrGroup(t *testing.T) {
	testsuite.Run(t, solutions.PromiseAllErrGroup[int])
}

func BenchmarkErrGroup(b *testing.B) {
	testsuite.Benchmark(b, solutions.PromiseAllErrGroup[int])
}
