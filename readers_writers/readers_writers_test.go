//go:build challenge

package readerswriters_test

import (
	"testing"

	readerswriters "github.com/N8Brooks/golang-concurrency/readers_writers"
	"github.com/N8Brooks/golang-concurrency/readers_writers/testsuite"
)

func TestReadersWriters(t *testing.T) {
	testsuite.Run(t, func() testsuite.ReadersWriters {
		return readerswriters.NewReadersWriters()
	})
}

func BenchmarkReadersWriters(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ReadersWriters {
		return readerswriters.NewReadersWriters()
	})
}
