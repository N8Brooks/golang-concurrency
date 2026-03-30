package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/cigarette_smokers/solutions"
	"github.com/N8Brooks/golang-concurrency/cigarette_smokers/testsuite"
)

func TestMutex(t *testing.T) {
	testsuite.Run(t, func(agent testsuite.Agent) testsuite.CigaretteSmokers {
		return solutions.NewMutex(agent)
	})
}

func BenchmarkMutex(b *testing.B) {
	testsuite.Benchmark(b, func(agent testsuite.Agent) testsuite.CigaretteSmokers {
		return solutions.NewMutex(agent)
	})
}
