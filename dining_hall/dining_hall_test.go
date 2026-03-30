//go:build challenge

package dininghall_test

import (
	"testing"

	dininghall "github.com/N8Brooks/golang-concurrency/dining_hall"
	"github.com/N8Brooks/golang-concurrency/dining_hall/testsuite"
)

func TestDiningHall(t *testing.T) {
	testsuite.Run(t, func() testsuite.DiningHall {
		return dininghall.NewDiningHall()
	})
}

func BenchmarkDiningHall(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.DiningHall {
		return dininghall.NewDiningHall()
	})
}
