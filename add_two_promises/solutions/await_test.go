package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/add_two_promises/solutions"
	"github.com/N8Brooks/golang-concurrency/add_two_promises/testsuite"
)

func TestAwait(t *testing.T) {
	testsuite.Run(t, solutions.AddTwoPromises)
}

func BenchmarkAwait(b *testing.B) {
	testsuite.Benchmark(b, solutions.AddTwoPromises)
}
