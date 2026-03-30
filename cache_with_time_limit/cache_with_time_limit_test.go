//go:build challenge

package cache_with_time_limit_test

import (
	"testing"

	cachewithtimelimit "github.com/N8Brooks/golang-concurrency/cache_with_time_limit"
	"github.com/N8Brooks/golang-concurrency/cache_with_time_limit/testsuite"
)

func TestTimeLimitedCache(t *testing.T) {
	testsuite.Run(t, func() testsuite.Cache {
		return cachewithtimelimit.NewTimeLimitedCache()
	})
}

func BenchmarkTimeLimitedCache(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.Cache {
		return cachewithtimelimit.NewTimeLimitedCache()
	})
}
