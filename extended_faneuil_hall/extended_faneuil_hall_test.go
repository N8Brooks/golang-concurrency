//go:build challenge

package extendedfaneuilhall_test

import (
	"testing"

	extendedfaneuilhall "github.com/N8Brooks/golang-concurrency/extended_faneuil_hall"
	"github.com/N8Brooks/golang-concurrency/extended_faneuil_hall/testsuite"
)

func TestExtendedFaneuilHall(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExtendedFaneuilHall {
		return extendedfaneuilhall.NewExtendedFaneuilHall()
	})
}

func BenchmarkExtendedFaneuilHall(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExtendedFaneuilHall {
		return extendedfaneuilhall.NewExtendedFaneuilHall()
	})
}
