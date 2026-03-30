//go:build challenge

package extendedchildcare_test

import (
	"testing"

	extendedchildcare "github.com/N8Brooks/golang-concurrency/extended_child_care"
	"github.com/N8Brooks/golang-concurrency/extended_child_care/testsuite"
)

func TestExtendedChildCare(t *testing.T) {
	testsuite.Run(t, func() testsuite.ExtendedChildCare {
		return extendedchildcare.NewExtendedChildCare()
	})
}

func BenchmarkExtendedChildCare(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ExtendedChildCare {
		return extendedchildcare.NewExtendedChildCare()
	})
}
