package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/generalized_smokers/solutions"
	"github.com/N8Brooks/golang-concurrency/generalized_smokers/testsuite"
)

func TestScoreboard(t *testing.T) {
	testsuite.Run(t, func(a testsuite.Agent) testsuite.Smokers {
		return solutions.NewScoreboard(a)
	})
}

func BenchmarkScoreboard(b *testing.B) {
	testsuite.Benchmark(b, func(a testsuite.Agent) testsuite.Smokers {
		return solutions.NewScoreboard(a)
	})
}
