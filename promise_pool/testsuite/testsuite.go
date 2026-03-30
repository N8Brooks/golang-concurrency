// Package testsuite contains reusable behavioral tests for Promise Pool implementations.
package testsuite

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type PromisePoolFunc func([]func() promise.Promiser[any], int) promise.Promiser[any]

func resolveAfter(d time.Duration) func() promise.Promiser[any] {
	return func() promise.Promiser[any] {
		return promise.New(func(resolve func(any), reject func(error)) {
			time.Sleep(d)
			resolve(struct{}{})
		})
	}
}

func Run(t *testing.T, promisePool PromisePoolFunc) {
	t.Helper()

	const threshold = 10 * time.Millisecond

	testCases := []struct {
		times []time.Duration
		n     int
		want  time.Duration
	}{
		{
			times: []time.Duration{300 * time.Millisecond, 400 * time.Millisecond, 200 * time.Millisecond},
			n:     2,
			want:  500 * time.Millisecond,
		},
		{
			times: []time.Duration{300 * time.Millisecond, 400 * time.Millisecond, 200 * time.Millisecond},
			n:     5,
			want:  400 * time.Millisecond,
		},
		{
			times: []time.Duration{300 * time.Millisecond, 400 * time.Millisecond, 200 * time.Millisecond},
			n:     1,
			want:  900 * time.Millisecond,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				functions := make([]func() promise.Promiser[any], len(tc.times))
				for i, d := range tc.times {
					functions[i] = resolveAfter(d)
				}

				start := time.Now()
				p := promisePool(functions, tc.n)
				if _, err := p.Result(); err != nil {
					t.Errorf("PromisePool(...) = _, %v, want _, nil", err)
				}
				synctest.Wait()
				duration := time.Since(start)

				if diff := (duration - tc.want).Abs(); diff > threshold {
					t.Errorf("PromisePool(...) took %v, want %v", duration, tc.want)
				}
			})
		})
	}
}

func Benchmark(b *testing.B, promisePool PromisePoolFunc) {
	b.Helper()
	b.ReportAllocs()

	functions := []func() promise.Promiser[any]{
		func() promise.Promiser[any] { return promise.Resolve(any(1)) },
		func() promise.Promiser[any] { return promise.Resolve(any(2)) },
		func() promise.Promiser[any] { return promise.Resolve(any(3)) },
	}

	for b.Loop() {
		if _, err := promisePool(functions, 2).Result(); err != nil {
			b.Fatalf("PromisePool() returned %v", err)
		}
	}
}
