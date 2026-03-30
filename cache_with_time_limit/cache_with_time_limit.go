//go:build challenge

// Package cache_with_time_limit contains the challenge version of the cache
// with TTL problem.
package cache_with_time_limit

import "time"

type TimeLimitedCache struct{}

func NewTimeLimitedCache() *TimeLimitedCache {
	return &TimeLimitedCache{}
}

func (c *TimeLimitedCache) Set(key int, value int, duration time.Duration) {
	panic("unimplemented")
}

func (c *TimeLimitedCache) Get(key int) int {
	panic("unimplemented")
}

func (c *TimeLimitedCache) Count() int {
	panic("unimplemented")
}
