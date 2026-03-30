//go:build challenge

package fifo_barbershop_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/fifo_barbershop"
	"github.com/N8Brooks/golang-concurrency/fifo_barbershop/testsuite"
)

func TestFIFOBarbershop(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.FIFOBarbershop {
		return fifo_barbershop.NewFIFOBarbershop(capacity)
	})
}

func BenchmarkFIFOBarbershop(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.FIFOBarbershop {
		return fifo_barbershop.NewFIFOBarbershop(capacity)
	})
}
