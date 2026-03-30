package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/mutex/solutions"
	"github.com/N8Brooks/golang-concurrency/mutex/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.Mutex {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Mutex {
		return solutions.NewChannel()
	})
}
