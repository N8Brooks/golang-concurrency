//go:build challenge

package multi_car_roller_coaster_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/multi_car_roller_coaster"
	"github.com/N8Brooks/golang-concurrency/multi_car_roller_coaster/testsuite"
)

func TestMultiCarRollerCoaster(t *testing.T) {
	testsuite.Run(t, func(cars, capacity int) testsuite.MultiCarRollerCoaster {
		return multi_car_roller_coaster.NewMultiCarRollerCoaster(cars, capacity)
	})
}

func BenchmarkMultiCarRollerCoaster(b *testing.B) {
	testsuite.Benchmark(b, func(cars, capacity int) testsuite.MultiCarRollerCoaster {
		return multi_car_roller_coaster.NewMultiCarRollerCoaster(cars, capacity)
	})
}
