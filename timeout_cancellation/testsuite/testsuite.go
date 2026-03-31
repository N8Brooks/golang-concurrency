// Package testsuite contains reusable behavioral tests for timeout cancellation implementations.
package testsuite

import (
	"reflect"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

type CancellableFunc func(fn func(...int), args []int, t time.Duration) func()

type invocation struct {
	at   time.Duration
	args []int
}

func Run(t *testing.T, cancellable CancellableFunc) {
	t.Helper()

	testCases := []struct {
		name     string
		args     []int
		delay    time.Duration
		cancelAt time.Duration
		expected []invocation
	}{
		{
			name:     "ExecutesBeforeCancellation",
			args:     []int{2},
			delay:    20 * time.Millisecond,
			cancelAt: 50 * time.Millisecond,
			expected: []invocation{{at: 20 * time.Millisecond, args: []int{2}}},
		},
		{
			name:     "CancelledBeforeExecution",
			args:     []int{3},
			delay:    50 * time.Millisecond,
			cancelAt: 20 * time.Millisecond,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

				cancel := cancellable(fn, tc.args, tc.delay)
				time.Sleep(tc.cancelAt)
				synctest.Wait()
				cancel()
				synctest.Wait()

				mu.Lock()
				defer mu.Unlock()
				if !reflect.DeepEqual(actual, tc.expected) {
					t.Fatalf("Cancellable() = %v, want %v", actual, tc.expected)
				}
			})
		})
	}
}
