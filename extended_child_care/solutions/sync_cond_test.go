package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/extended_child_care/solutions"
	"github.com/N8Brooks/golang-concurrency/extended_child_care/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExtendedChildCare {
		return solutions.NewSyncCond()
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExtendedChildCare {
		return solutions.NewSyncCond()
	})
}
