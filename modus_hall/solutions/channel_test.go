package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/modus_hall/solutions"
	"github.com/N8Brooks/golang-concurrency/modus_hall/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.ModusHall {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ModusHall {
		return solutions.NewChannel()
	})
}
