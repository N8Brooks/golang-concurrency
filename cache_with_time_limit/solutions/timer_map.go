// Package solutions contains implementations of the cache with time limit problem.
package solutions

import (
	"sync"
	"time"
)

type TimeLimitedCache struct {
	cache map[int]entry
	mu    sync.RWMutex
}

type entry struct {
	gen   uint64
	value int
	timer *time.Timer
}

func NewTimeLimitedCache() *TimeLimitedCache {
	return &TimeLimitedCache{cache: make(map[int]entry)}
}

func (c *TimeLimitedCache) Set(key int, value int, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var gen uint64 = 1
	if e, ok := c.cache[key]; ok {
		e.timer.Stop()
		gen = e.gen + 1
	}

	timer := time.AfterFunc(duration, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		if current, ok := c.cache[key]; ok && current.gen == gen {
			delete(c.cache, key)
		}
	})

	c.cache[key] = entry{gen: gen, value: value, timer: timer}
}

func (c *TimeLimitedCache) Get(key int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.cache[key]; ok {
		return e.value
	}
	return -1
}

func (c *TimeLimitedCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}
