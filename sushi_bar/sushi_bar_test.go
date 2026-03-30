//go:build challenge

package sushibar_test

import (
	"testing"

	sushibar "github.com/N8Brooks/golang-concurrency/sushi_bar"
	"github.com/N8Brooks/golang-concurrency/sushi_bar/testsuite"
)

func TestSushiBar(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.SushiBar {
		return sushibar.NewSushiBar(capacity)
	})
}

func BenchmarkSushiBar(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.SushiBar {
		return sushibar.NewSushiBar(capacity)
	})
}
