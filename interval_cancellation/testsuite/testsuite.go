// Package testsuite contains reusable behavioral tests for interval cancellation implementations.
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

		cancel := cancellable(fn, []int{4}, 35*time.Millisecond)
		time.Sleep(190 * time.Millisecond)
		synctest.Wait()
		cancel()
		synctest.Wait()

		mu.Lock()
		defer mu.Unlock()
		want := []invocation{
			{at: 0, args: []int{4}},
			{at: 35 * time.Millisecond, args: []int{4}},
			{at: 70 * time.Millisecond, args: []int{4}},
			{at: 105 * time.Millisecond, args: []int{4}},
			{at: 140 * time.Millisecond, args: []int{4}},
			{at: 175 * time.Millisecond, args: []int{4}},
		}
		if !reflect.DeepEqual(actual, want) {
			t.Fatalf("Cancellable() = %v, want %v", actual, want)
		}
	})
}
