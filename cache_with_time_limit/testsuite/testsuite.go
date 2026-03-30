// Package testsuite contains reusable behavioral tests for cache with time limit implementations.
package testsuite

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"
)

type Cache interface {
	Set(key int, value int, duration time.Duration)
	Get(key int) int
	Count() int
}

type actor interface {
	act(c Cache, t *testing.T)
}

type set struct {
	key      int
	value    int
	duration time.Duration
}

func (s *set) act(c Cache, t *testing.T) {
	c.Set(s.key, s.value, s.duration)
}

type get struct {
	key   int
	value int
}

func (g *get) act(c Cache, t *testing.T) {
	if v := c.Get(g.key); v != g.value {
		t.Errorf("Get(%d) = %d, want %d", g.key, v, g.value)
	}
}

type count struct {
	want int
}

func (cnt *count) act(c Cache, t *testing.T) {
	if v := c.Count(); v != cnt.want {
		t.Errorf("Count() = %d, want %d", v, cnt.want)
	}
}

func Run(t *testing.T, newImpl func() Cache) {
	t.Helper()

	testCases := []struct {
		actions []actor
		at      []time.Duration
	}{
		{
			actions: []actor{
				&set{1, 42, 100 * time.Millisecond},
				&get{1, 42},
				&count{1},
				&get{1, -1},
			},
			at: []time.Duration{
				0,
				50 * time.Millisecond,
				50 * time.Millisecond,
				150 * time.Millisecond,
			},
		},
		{
			actions: []actor{
				&set{1, 42, 50 * time.Millisecond},
				&set{1, 50, 100 * time.Millisecond},
				&get{1, 50},
				&get{1, 50},
				&get{1, -1},
				&count{0},
			},
			at: []time.Duration{
				0,
				40 * time.Millisecond,
				50 * time.Millisecond,
				120 * time.Millisecond,
				200 * time.Millisecond,
				250 * time.Millisecond,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				var current time.Duration
				c := newImpl()

				for i, action := range tc.actions {
					next := tc.at[i]
					time.Sleep(next - current)
					synctest.Wait()
					current = next
					action.act(c, t)
				}
			})
		})
	}
}

func Benchmark(b *testing.B, newImpl func() Cache) {
	b.Helper()
	b.ReportAllocs()

	c := newImpl()
	for i := 0; b.Loop(); i++ {
		c.Set(i, i, time.Hour)
		if got := c.Get(i); got != i {
			b.Fatalf("Get(%d) = %d", i, got)
		}
	}
}
