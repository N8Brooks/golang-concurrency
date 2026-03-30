//go:build challenge

package writerspriorityreaderswriters_test

import (
	"testing"

	writerspriorityreaderswriters "github.com/N8Brooks/golang-concurrency/writers_priority_readers_writers"
	"github.com/N8Brooks/golang-concurrency/writers_priority_readers_writers/testsuite"
)

func TestWritersPriorityReadersWriters(t *testing.T) {
	testsuite.Run(t, func() testsuite.WritersPriorityReadersWriters {
		return writerspriorityreaderswriters.NewWritersPriorityReadersWriters()
	})
}

func BenchmarkWritersPriorityReadersWriters(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.WritersPriorityReadersWriters {
		return writerspriorityreaderswriters.NewWritersPriorityReadersWriters()
	})
}
