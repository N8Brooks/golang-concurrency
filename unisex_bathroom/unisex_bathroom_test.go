//go:build challenge

package unisexbathroom_test

import (
	"testing"

	unisexbathroom "github.com/N8Brooks/golang-concurrency/unisex_bathroom"
	"github.com/N8Brooks/golang-concurrency/unisex_bathroom/testsuite"
)

func TestUnisexBathroom(t *testing.T) {
	testsuite.Run(t, func() testsuite.UnisexBathroom {
		return unisexbathroom.NewUnisexBathroom()
	})
}

func BenchmarkUnisexBathroom(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.UnisexBathroom {
		return unisexbathroom.NewUnisexBathroom()
	})
}
