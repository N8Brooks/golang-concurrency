package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/hilzers_barbershop/solutions"
	"github.com/N8Brooks/golang-concurrency/hilzers_barbershop/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.Barbershop {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Barbershop {
		return solutions.NewSyncMutex()
	})
}
