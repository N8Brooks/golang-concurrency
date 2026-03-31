// Package testsuite contains reusable behavioral tests for query batching implementations.
package testsuite

import (
	"fmt"
	"reflect"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type Batcher interface {
	GetValue(key string) promise.Promiser[string]
}

func Run(t *testing.T, newImpl func(queryFn func([]string) promise.Promiser[[]string], d time.Duration) Batcher) {
	t.Helper()

	synctest.Test(t, func(t *testing.T) {
		var mu sync.Mutex
		var batches [][]string

		queryFn := func(keys []string) promise.Promiser[[]string] {
			mu.Lock()
			batches = append(batches, append([]string(nil), keys...))
			mu.Unlock()

			return promise.New(func(resolve func([]string), reject func(error)) {
				time.Sleep(10 * time.Millisecond)
				results := make([]string, len(keys))
				for i, key := range keys {
					results[i] = fmt.Sprintf("value:%s", key)
				}
				resolve(results)
			})
		}

		batcher := newImpl(queryFn, 100*time.Millisecond)
		var current time.Duration

		p1 := batcher.GetValue("a")
		time.Sleep(10*time.Millisecond - current)
		synctest.Wait()
		current = 10 * time.Millisecond

		p2 := batcher.GetValue("b")
		time.Sleep(10 * time.Millisecond)
		synctest.Wait()
		current = 20 * time.Millisecond

		p3 := batcher.GetValue("c")
		time.Sleep(10 * time.Millisecond)
		synctest.Wait()
		current = 30 * time.Millisecond

		p4 := batcher.GetValue("d")
		time.Sleep(220*time.Millisecond - current)
		synctest.Wait()
		current = 220 * time.Millisecond

		p5 := batcher.GetValue("e")
		time.Sleep(50 * time.Millisecond)
		synctest.Wait()

		wantValues := []string{"value:a", "value:b", "value:c", "value:d", "value:e"}
		promises := []promise.Promiser[string]{p1, p2, p3, p4, p5}
		for i, p := range promises {
			if got, err := p.Result(); err != nil || got != wantValues[i] {
				t.Fatalf("GetValue() = %q, %v, want %q, nil", got, err, wantValues[i])
			}
		}

		mu.Lock()
		defer mu.Unlock()
		wantBatches := [][]string{
			{"a"},
			{"b", "c", "d"},
			{"e"},
		}
		if !reflect.DeepEqual(batches, wantBatches) {
			t.Fatalf("query batches = %v, want %v", batches, wantBatches)
		}

		// Ensure batching preserves input order.
		if !slices.Equal(batches[1], []string{"b", "c", "d"}) {
			t.Fatalf("batched keys = %v, want [b c d]", batches[1])
		}
	})
}
