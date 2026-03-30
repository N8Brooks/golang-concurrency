package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/writers_priority_readers_writers/solutions"
	"github.com/N8Brooks/golang-concurrency/writers_priority_readers_writers/testsuite"
)

func TestSyncCond(t *testing.T) {
	testsuite.Run(t, func() testsuite.WritersPriorityReadersWriters {
		return solutions.NewSyncCond()
	})
}

func BenchmarkSyncCond(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.WritersPriorityReadersWriters {
		return solutions.NewSyncCond()
	})
}
