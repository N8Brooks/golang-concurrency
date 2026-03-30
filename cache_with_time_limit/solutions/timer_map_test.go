package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/cache_with_time_limit/solutions"
	"github.com/N8Brooks/golang-concurrency/cache_with_time_limit/testsuite"
)

func TestTimerMap(t *testing.T) {
	testsuite.Run(t, func() testsuite.Cache {
		return solutions.NewTimeLimitedCache()
	})
}

func BenchmarkTimerMap(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Cache {
		return solutions.NewTimeLimitedCache()
	})
}
