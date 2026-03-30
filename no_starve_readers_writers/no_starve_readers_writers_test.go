//go:build challenge

package nostarvereaderswriters_test

import (
	"testing"

	nostarvereaderswriters "github.com/N8Brooks/golang-concurrency/no_starve_readers_writers"
	"github.com/N8Brooks/golang-concurrency/no_starve_readers_writers/testsuite"
)

func TestNoStarveReadersWriters(t *testing.T) {
	testsuite.Run(t, func() testsuite.NoStarveReadersWriters {
		return nostarvereaderswriters.NewNoStarveReadersWriters()
	})
}

func BenchmarkNoStarveReadersWriters(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.NoStarveReadersWriters {
		return nostarvereaderswriters.NewNoStarveReadersWriters()
	})
}
