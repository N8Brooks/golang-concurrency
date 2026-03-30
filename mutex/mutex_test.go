//go:build challenge

package mutex_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/mutex"
	"github.com/N8Brooks/golang-concurrency/mutex/testsuite"
)

func TestMutex(t *testing.T) {
	testsuite.Run(t, func() testsuite.Mutex {
		return mutex.NewMutex()
	})
}

func BenchmarkMutex(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Mutex {
		return mutex.NewMutex()
	})
}
