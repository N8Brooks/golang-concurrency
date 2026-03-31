// Package testsuite contains reusable behavioral tests for custom interval implementations.
package testsuite

import (
	"reflect"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

type CustomIntervalFunc func(fn func(), delay, period time.Duration) int
type ClearIntervalFunc func(id int)

func Run(t *testing.T, set CustomIntervalFunc, clear ClearIntervalFunc) {
	t.Helper()

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		var mu sync.Mutex
		var actual []time.Duration

		id := set(func() {
			mu.Lock()
			defer mu.Unlock()
			actual = append(actual, time.Since(start))
		}, 50*time.Millisecond, 20*time.Millisecond)

		time.Sleep(225 * time.Millisecond)
		synctest.Wait()
		clear(id)
		synctest.Wait()

		mu.Lock()
		defer mu.Unlock()
		want := []time.Duration{
			50 * time.Millisecond,
			120 * time.Millisecond,
			210 * time.Millisecond,
		}
		if !reflect.DeepEqual(actual, want) {
			t.Fatalf("CustomInterval() = %v, want %v", actual, want)
		}
	})
}
