//go:build challenge

package barbershop_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/barbershop"
	"github.com/N8Brooks/golang-concurrency/barbershop/testsuite"
)

func TestBarbershop(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.Barbershop {
		return barbershop.NewBarbershop(capacity)
	})
}

func BenchmarkBarbershop(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.Barbershop {
		return barbershop.NewBarbershop(capacity)
	})
}
