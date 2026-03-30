//go:build challenge

package senatebus_test

import (
	"testing"

	senatebus "github.com/N8Brooks/golang-concurrency/senate_bus"
	"github.com/N8Brooks/golang-concurrency/senate_bus/testsuite"
)

func TestSenateBus(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.SenateBus {
		return senatebus.NewSenateBus(capacity)
	})
}

func BenchmarkSenateBus(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.SenateBus {
		return senatebus.NewSenateBus(capacity)
	})
}
