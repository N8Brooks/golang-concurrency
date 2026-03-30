package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/child_care/solutions"
	"github.com/N8Brooks/golang-concurrency/child_care/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func() testsuite.ChildCare {
		return solutions.NewSyncCond()
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ChildCare {
		return solutions.NewSyncCond()
	})
}
