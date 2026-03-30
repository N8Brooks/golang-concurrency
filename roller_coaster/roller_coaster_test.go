//go:build challenge

package roller_coaster_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/roller_coaster"
	"github.com/N8Brooks/golang-concurrency/roller_coaster/testsuite"
)

func TestRollerCoaster(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.RollerCoaster {
		return roller_coaster.NewRollerCoaster(capacity)
	})
}

func BenchmarkRollerCoaster(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.RollerCoaster {
		return roller_coaster.NewRollerCoaster(capacity)
	})
}
