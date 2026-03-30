//go:build challenge

package cigarettesmokers_test

import (
	"testing"

	cigarettesmokers "github.com/N8Brooks/golang-concurrency/cigarette_smokers"
	"github.com/N8Brooks/golang-concurrency/cigarette_smokers/testsuite"
)

func TestCigaretteSmokers(t *testing.T) {
	testsuite.Run(t, func(agent testsuite.Agent) testsuite.CigaretteSmokers {
		return cigarettesmokers.NewCigaretteSmokers(agent)
	})
}

func BenchmarkCigaretteSmokers(b *testing.B) {
	testsuite.Benchmark(b, func(agent testsuite.Agent) testsuite.CigaretteSmokers {
		return cigarettesmokers.NewCigaretteSmokers(agent)
	})
}
