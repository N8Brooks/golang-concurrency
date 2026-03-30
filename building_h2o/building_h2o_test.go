//go:build challenge

package buildingh2o_test

import (
	"testing"

	buildingh2o "github.com/N8Brooks/golang-concurrency/building_h2o"
	"github.com/N8Brooks/golang-concurrency/building_h2o/testsuite"
)

func TestH2O(t *testing.T) {
	testsuite.Run(t, func() testsuite.H2O {
		return buildingh2o.NewH2O()
	})
}

func BenchmarkH2O(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.H2O {
		return buildingh2o.NewH2O()
	})
}
