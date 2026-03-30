//go:build challenge

package print_in_order_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/print_in_order"
	"github.com/N8Brooks/golang-concurrency/print_in_order/testsuite"
)

func TestPrintInOrder(t *testing.T) {
	testsuite.Run(t, func() testsuite.PrintInOrder {
		return print_in_order.NewPrintInOrder()
	})
}

func BenchmarkPrintInOrder(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.PrintInOrder {
		return print_in_order.NewPrintInOrder()
	})
}
