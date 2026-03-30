// Package testsuite contains reusable behavioral tests for Promise Time Limit implementations.
package testsuite

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	promisetimelimit "github.com/N8Brooks/golang-concurrency/promise_time_limit"
)

type TimeLimitFunc func(fn func(...int) promise.Promiser[int], t time.Duration) func(...int) promise.Promiser[int]

var Err = errors.New("error")

func Run(t *testing.T, timeLimit TimeLimitFunc) {
	t.Helper()

	const threshold = 10 * time.Millisecond

	testCases := []struct {
		fn     func(...int) promise.Promiser[int]
		inputs []int
		limit  time.Duration
		res    int
		err    error
	}{
		{
			fn: func(args ...int) promise.Promiser[int] {
				time.Sleep(100 * time.Millisecond)
				return promise.Resolve(args[0] * args[0])
			},
			inputs: []int{5},
			limit:  50 * time.Millisecond,
			err:    promisetimelimit.ErrTimeLimitExceeded,
		},
		{
			fn: func(args ...int) promise.Promiser[int] {
				time.Sleep(100 * time.Millisecond)
				return promise.Resolve(args[0] * args[0])
			},
			inputs: []int{5},
			limit:  150 * time.Millisecond,
			res:    25,
		},
		{
			fn: func(args ...int) promise.Promiser[int] {
				time.Sleep(120 * time.Millisecond)
				return promise.Resolve(args[0] + args[1])
			},
			inputs: []int{5, 10},
			limit:  150 * time.Millisecond,
			res:    15,
		},
		{
			fn: func(args ...int) promise.Promiser[int] {
				return promise.Reject[int](Err)
			},
			limit: 1000 * time.Millisecond,
			err:   Err,
		},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			fn := timeLimit(tt.fn, tt.limit)
			start := time.Now()
			result, err := fn(tt.inputs...).Result()
			duration := time.Since(start)
			if duration > tt.limit+threshold {
				t.Errorf("TimeLimit() took %v, want <= %v", duration, tt.limit+threshold)
			}
			if result != tt.res || !errors.Is(err, tt.err) {
				t.Errorf("TimeLimit() = %d, %v, want %d, %v", result, err, tt.res, tt.err)
			}
		})
	}
}

func Benchmark(b *testing.B, timeLimit TimeLimitFunc) {
	b.Helper()
	b.ReportAllocs()

	limited := timeLimit(func(args ...int) promise.Promiser[int] {
		return promise.Resolve(args[0] + 1)
	}, time.Second)

	for i := 0; b.Loop(); i++ {
		if got, err := limited(i).Result(); err != nil || got != i+1 {
			b.Fatalf("TimeLimit() = %d, %v", got, err)
		}
	}
}
