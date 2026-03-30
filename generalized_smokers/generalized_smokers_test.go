//go:build challenge

package generalizedsmokers_test

import (
	"testing"

	generalizedsmokers "github.com/N8Brooks/golang-concurrency/generalized_smokers"
	"github.com/N8Brooks/golang-concurrency/generalized_smokers/testsuite"
)

func TestGeneralizedSmokers(t *testing.T) {
	testsuite.Run(t, func(a testsuite.Agent) testsuite.Smokers {
		return generalizedsmokers.NewSmokers(a)
	})
}

func BenchmarkGeneralizedSmokers(b *testing.B) {
	testsuite.Benchmark(b, func(a testsuite.Agent) testsuite.Smokers {
		return generalizedsmokers.NewSmokers(a)
	})
}
