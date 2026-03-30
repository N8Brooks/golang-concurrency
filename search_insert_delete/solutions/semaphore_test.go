package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/search_insert_delete/solutions"
	"github.com/N8Brooks/golang-concurrency/search_insert_delete/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.SearchInsertDelete {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.SearchInsertDelete {
		return solutions.NewSemaphore()
	})
}
