//go:build challenge

package rendezvous_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/rendezvous"
	"github.com/N8Brooks/golang-concurrency/rendezvous/testsuite"
)

func TestRendezvous(t *testing.T) {
	testsuite.Run(t, func() testsuite.Rendezvous {
		return rendezvous.NewRendezvous()
	})
}

func BenchmarkRendezvous(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Rendezvous {
		return rendezvous.NewRendezvous()
	})
}
