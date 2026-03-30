//go:build challenge

package hilzersbarbershop_test

import (
	"testing"

	hilzersbarbershop "github.com/N8Brooks/golang-concurrency/hilzers_barbershop"
	"github.com/N8Brooks/golang-concurrency/hilzers_barbershop/testsuite"
)

func TestHilzersBarbershop(t *testing.T) {
	testsuite.Run(t, func() testsuite.Barbershop {
		return hilzersbarbershop.NewBarbershop()
	})
}

func BenchmarkHilzersBarbershop(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Barbershop {
		return hilzersbarbershop.NewBarbershop()
	})
}
