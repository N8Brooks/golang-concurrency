//go:build challenge

package add_two_promises_test

import (
	"testing"

	addtwopromises "github.com/N8Brooks/golang-concurrency/add_two_promises"
	"github.com/N8Brooks/golang-concurrency/add_two_promises/testsuite"
)

func TestAddTwoPromises(t *testing.T) {
	testsuite.Run(t, addtwopromises.AddTwoPromises)
}

func BenchmarkAddTwoPromises(b *testing.B) {
	testsuite.Benchmark(b, addtwopromises.AddTwoPromises)
}
