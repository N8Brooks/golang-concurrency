//go:build challenge

package santaclaus_test

import (
	"testing"

	santaclaus "github.com/N8Brooks/golang-concurrency/santa_claus"
	"github.com/N8Brooks/golang-concurrency/santa_claus/testsuite"
)

func TestSantaClaus(t *testing.T) {
	testsuite.Run(t, func() testsuite.SantaClaus {
		return santaclaus.NewSantaClaus()
	})
}

func BenchmarkSantaClaus(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.SantaClaus {
		return santaclaus.NewSantaClaus()
	})
}
