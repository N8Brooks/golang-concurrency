//go:build challenge

package babooncrossing_test

import (
	"testing"

	babooncrossing "github.com/N8Brooks/golang-concurrency/baboon_crossing"
	"github.com/N8Brooks/golang-concurrency/baboon_crossing/testsuite"
)

func TestBaboonCrossing(t *testing.T) {
	testsuite.Run(t, func() testsuite.Crossing {
		return babooncrossing.NewCrossing()
	})
}

func BenchmarkBaboonCrossing(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Crossing {
		return babooncrossing.NewCrossing()
	})
}
