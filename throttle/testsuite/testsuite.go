// Package testsuite contains reusable behavioral tests for throttle implementations.
package testsuite

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

type ThrottleFunc func(fn func(...int), d time.Duration) func(...int)

type call struct {
	at   time.Duration
	args []int
}

type invocation struct {
	at   time.Duration
	args []int
}

func Run(t *testing.T, throttle ThrottleFunc) {
	t.Helper()

	testCases := []struct {
		delay    time.Duration
		calls    []call
		finishAt time.Duration
		expected []invocation
	}{
		{
			delay: 50 * time.Millisecond,
			calls: []call{
				{at: 0, args: []int{1}},
				{at: 10 * time.Millisecond, args: []int{2}},
				{at: 20 * time.Millisecond, args: []int{3}},
			},
			finishAt: 50 * time.Millisecond,
			expected: []invocation{
				{at: 0, args: []int{1}},
				{at: 50 * time.Millisecond, args: []int{3}},
			},
		},
		{
			delay: 50 * time.Millisecond,
			calls: []call{
				{at: 0, args: []int{1}},
				{at: 60 * time.Millisecond, args: []int{2}},
			},
			finishAt: 60 * time.Millisecond,
			expected: []invocation{
				{at: 0, args: []int{1}},
				{at: 60 * time.Millisecond, args: []int{2}},
			},
		},
		{
			delay: 50 * time.Millisecond,
			calls: []call{
				{at: 0, args: []int{1}},
				{at: 10 * time.Millisecond, args: []int{2}},
				{at: 20 * time.Millisecond, args: []int{3}},
				{at: 70 * time.Millisecond, args: []int{4}},
			},
			finishAt: 100 * time.Millisecond,
			expected: []invocation{
				{at: 0, args: []int{1}},
				{at: 50 * time.Millisecond, args: []int{3}},
				{at: 100 * time.Millisecond, args: []int{4}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()
				var mu sync.Mutex
				var actual []invocation

				fn := func(args ...int) {
					mu.Lock()
					defer mu.Unlock()
					actual = append(actual, invocation{
						at:   time.Since(start),
						args: append([]int(nil), args...),
					})
				}

				throttled := throttle(fn, tc.delay)
				var current time.Duration

				for _, call := range tc.calls {
					time.Sleep(call.at - current)
					synctest.Wait()
					current = call.at
					throttled(call.args...)
				}

				if tc.finishAt > current {
					time.Sleep(tc.finishAt - current)
				}
				synctest.Wait()

				mu.Lock()
				defer mu.Unlock()
				if !reflect.DeepEqual(actual, tc.expected) {
					t.Fatalf("Throttle() = %v, want %v", actual, tc.expected)
				}
			})
		})
	}
}
