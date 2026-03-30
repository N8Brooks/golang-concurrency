//go:build challenge

package no_starve_mutex_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/no_starve_mutex"
	"github.com/N8Brooks/golang-concurrency/no_starve_mutex/testsuite"
)

func TestNoStarveMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.Mutex {
		return no_starve_mutex.NewNoStarveMutex()
	})
}

func BenchmarkNoStarveMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Mutex {
		return no_starve_mutex.NewNoStarveMutex()
	})
}
