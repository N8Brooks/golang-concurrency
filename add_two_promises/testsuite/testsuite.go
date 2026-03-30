// Package testsuite contains reusable behavioral tests for Add Two Promises implementations.
package testsuite

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type AddTwoPromisesFunc func(a, b promise.Promiser[int]) promise.Promiser[int]

func intAfter(n int, d time.Duration) promise.Promiser[int] {
	return promise.New(func(resolve func(int), reject func(error)) {
		time.Sleep(d)
		resolve(n)
	})
}

func Run(t *testing.T, add AddTwoPromisesFunc) {
	t.Helper()

	testCases := []struct {
		x1, x2 int
		t1, t2 time.Duration
		want   int
	}{
		{x1: 2, x2: 5, t1: 20 * time.Millisecond, t2: 60 * time.Millisecond, want: 7},
		{x1: 10, x2: -12, t1: 50 * time.Millisecond, t2: 30 * time.Millisecond, want: -2},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				a := intAfter(tc.x1, tc.t1)
				b := intAfter(tc.x2, tc.t2)
				p := add(a, b)
				if res, err := p.Result(); res != tc.want || err != nil {
					t.Errorf("AddTwoPromises(...) = %d, %v, want %d, nil", res, err, tc.want)
				}
			})
		})
	}
}

func Benchmark(b *testing.B, add AddTwoPromisesFunc) {
	b.Helper()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		p := add(promise.Resolve(i), promise.Resolve(i+1))
		if got, err := p.Result(); err != nil || got != 2*i+1 {
			b.Fatalf("AddTwoPromises(...) = %d, %v", got, err)
		}
	}
}
