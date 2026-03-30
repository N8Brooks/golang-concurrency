//go:build challenge

package searchinsertdelete_test

import (
	"testing"

	searchinsertdelete "github.com/N8Brooks/golang-concurrency/search_insert_delete"
	"github.com/N8Brooks/golang-concurrency/search_insert_delete/testsuite"
)

func TestSearchInsertDelete(t *testing.T) {
	testsuite.Run(t, func() testsuite.SearchInsertDelete {
		return searchinsertdelete.NewSearchInsertDelete()
	})
}

func BenchmarkSearchInsertDelete(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.SearchInsertDelete {
		return searchinsertdelete.NewSearchInsertDelete()
	})
}
