package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/search_insert_delete/solutions"
	"github.com/N8Brooks/golang-concurrency/search_insert_delete/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.SearchInsertDelete {
		return solutions.NewSyncMutex()
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.SearchInsertDelete {
		return solutions.NewSyncMutex()
	})
}
