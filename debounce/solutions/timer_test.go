package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/debounce/solutions"
	"github.com/N8Brooks/golang-concurrency/debounce/testsuite"
)

func TestTimer(t *testing.T) {
	testsuite.Run(t, solutions.Debounce)
}

func BenchmarkTimer(b *testing.B) {
	testsuite.Benchmark(b, solutions.Debounce)
}
