//go:build challenge

package trafficlight_test

import (
	"testing"

	trafficlight "github.com/N8Brooks/golang-concurrency/traffic_light"
	"github.com/N8Brooks/golang-concurrency/traffic_light/testsuite"
)

func TestTrafficLight(t *testing.T) {
	testsuite.Run(t, func() testsuite.TrafficLight {
		return trafficlight.NewTrafficLight()
	})
}

func BenchmarkTrafficLight(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.TrafficLight {
		return trafficlight.NewTrafficLight()
	})
}
