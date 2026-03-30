package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/traffic_light/solutions"
	"github.com/N8Brooks/golang-concurrency/traffic_light/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.TrafficLight {
		return solutions.NewChannel()
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.TrafficLight {
		return solutions.NewChannel()
	})
}
