//go:build challenge

package faneuilhall_test

import (
	"testing"

	faneuilhall "github.com/N8Brooks/golang-concurrency/faneuil_hall"
	"github.com/N8Brooks/golang-concurrency/faneuil_hall/testsuite"
)

func TestFaneuilHall(t *testing.T) {
	testsuite.Run(t, func() testsuite.FaneuilHall {
		return faneuilhall.NewFaneuilHall()
	})
}

func BenchmarkFaneuilHall(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.FaneuilHall {
		return faneuilhall.NewFaneuilHall()
	})
}
