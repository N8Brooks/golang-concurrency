package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/sleep/solutions"
	"github.com/N8Brooks/golang-concurrency/sleep/testsuite"
)

func TestPromise(t *testing.T) {
	testsuite.Run(t, solutions.Sleep)
}

func BenchmarkPromise(b *testing.B) {
	testsuite.Benchmark(b, solutions.Sleep)
}
