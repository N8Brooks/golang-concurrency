//go:build challenge

package childcare_test

import (
	"testing"

	childcare "github.com/N8Brooks/golang-concurrency/child_care"
	"github.com/N8Brooks/golang-concurrency/child_care/testsuite"
)

func TestChildCare(t *testing.T) {
	testsuite.Run(t, func() testsuite.ChildCare {
		return childcare.NewChildCare()
	})
}

func BenchmarkChildCare(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ChildCare {
		return childcare.NewChildCare()
	})
}
