package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/multi_car_roller_coaster/solutions"
	"github.com/N8Brooks/golang-concurrency/multi_car_roller_coaster/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func(cars, capacity int) testsuite.MultiCarRollerCoaster {
		return solutions.NewSemaphore(cars, capacity)
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func(cars, capacity int) testsuite.MultiCarRollerCoaster {
		return solutions.NewSemaphore(cars, capacity)
	})
}
