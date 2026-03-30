//go:build challenge

package extendeddininghall_test

import (
	"testing"

	extendeddininghall "github.com/N8Brooks/golang-concurrency/extended_dining_hall"
	"github.com/N8Brooks/golang-concurrency/extended_dining_hall/testsuite"
)

func TestExtendedDiningHall(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExtendedDiningHall {
		return extendeddininghall.NewExtendedDiningHall()
	})
}

func BenchmarkExtendedDiningHall(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExtendedDiningHall {
		return extendeddininghall.NewExtendedDiningHall()
	})
}
