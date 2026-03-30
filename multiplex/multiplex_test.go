//go:build challenge

package multiplex_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/multiplex"
	"github.com/N8Brooks/golang-concurrency/multiplex/testsuite"
)

func TestMultiplex(t *testing.T) {
	testsuite.Run(t, func(limit int) testsuite.Multiplex {
		return multiplex.NewMultiplex(limit)
	})
}

func BenchmarkMultiplex(b *testing.B) {
	testsuite.Benchmark(b, func(limit int) testsuite.Multiplex {
		return multiplex.NewMultiplex(limit)
	})
}
