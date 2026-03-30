// Package testsuite contains reusable behavioral tests for debounce implementations.
package testsuite

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

type DebounceFunc func(fn func(...int), t time.Duration) func(...int)

type call struct {
	at     time.Duration
	inputs []int
}

func Run(t *testing.T, debounce DebounceFunc) {
	t.Helper()

	testCases := []struct {
		delay    time.Duration
		calls    []call
		expected [][]int
	}{
		{
			delay: 50 * time.Millisecond,
			calls: []call{
				{at: 50 * time.Millisecond, inputs: []int{1}},
				{at: 75 * time.Millisecond, inputs: []int{2}},
			},
			expected: [][]int{{2}},
		},
		{
			delay: 20 * time.Millisecond,
			calls: []call{
				{at: 50 * time.Millisecond, inputs: []int{1}},
				{at: 100 * time.Millisecond, inputs: []int{2}},
			},
			expected: [][]int{{1}, {2}},
		},
		{
			delay: 150 * time.Millisecond,
			calls: []call{
				{at: 50 * time.Millisecond, inputs: []int{1, 2}},
				{at: 300 * time.Millisecond, inputs: []int{3, 4}},
				{at: 300 * time.Millisecond, inputs: []int{5, 6}},
			},
			expected: [][]int{{1, 2}, {5, 6}},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				var mu sync.Mutex
				var actual [][]int

				fn := func(args ...int) {
					mu.Lock()
					defer mu.Unlock()
					actual = append(actual, append([]int(nil), args...))
				}

				debounced := debounce(fn, tc.delay)

				var current time.Duration
				for _, c := range tc.calls {
					time.Sleep(c.at - current)
					synctest.Wait()
					current = c.at
					debounced(c.inputs...)
				}

				time.Sleep(24 * time.Hour)
				synctest.Wait()

				if !reflect.DeepEqual(actual, tc.expected) {
					t.Errorf("Debounce() = %v, want %v", actual, tc.expected)
				}
			})
		})
	}
}

func Benchmark(b *testing.B, debounce DebounceFunc) {
	b.Helper()
	b.ReportAllocs()

	var mu sync.Mutex
	var total int
	debounced := debounce(func(args ...int) {
		mu.Lock()
		total += len(args)
		mu.Unlock()
	}, time.Nanosecond)

	for i := 0; b.Loop(); i++ {
		debounced(i)
	}

	_ = total
}
