//go:build challenge

package modushall_test

import (
	"testing"

	modushall "github.com/N8Brooks/golang-concurrency/modus_hall"
	"github.com/N8Brooks/golang-concurrency/modus_hall/testsuite"
)

func TestModusHall(t *testing.T) {
	testsuite.Run(t, func() testsuite.ModusHall {
		return modushall.NewModusHall()
	})
}

func BenchmarkModusHall(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.ModusHall {
		return modushall.NewModusHall()
	})
}
