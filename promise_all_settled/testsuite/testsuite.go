// Package testsuite contains reusable behavioral tests for promise_all_settled implementations.
package testsuite

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	promiseallsettled "github.com/N8Brooks/golang-concurrency/promise_all_settled"
)

type PromiseAllSettledFunc func([]func() promise.Promiser[int]) []promiseallsettled.Result[int]

func resolveAfter(n int, d time.Duration) func() promise.Promiser[int] {
	return func() promise.Promiser[int] {
		return promise.New(func(resolve func(int), reject func(error)) {
			time.Sleep(d)
			resolve(n)
		})
	}
}

func rejectAfter(err error, d time.Duration) func() promise.Promiser[int] {
	return func() promise.Promiser[int] {
		return promise.New(func(resolve func(int), reject func(error)) {
			time.Sleep(d)
			reject(err)
		})
	}
}

func Run(t *testing.T, promiseAllSettled PromiseAllSettledFunc) {
	t.Helper()

	const threshold = 10 * time.Millisecond
	errBoom := errors.New("error")

	tests := []struct {
		functions []func() promise.Promiser[int]
		want      []promiseallsettled.Result[int]
		elapsed   time.Duration
	}{
		{
			functions: []func() promise.Promiser[int]{
				resolveAfter(15, 100*time.Millisecond),
			},
			want: []promiseallsettled.Result[int]{
				{Val: 15},
			},
			elapsed: 100 * time.Millisecond,
		},
		{
			functions: []func() promise.Promiser[int]{
				resolveAfter(20, 100*time.Millisecond),
				resolveAfter(15, 100*time.Millisecond),
			},
			want: []promiseallsettled.Result[int]{
				{Val: 20},
				{Val: 15},
			},
			elapsed: 100 * time.Millisecond,
		},
		{
			functions: []func() promise.Promiser[int]{
				resolveAfter(30, 200*time.Millisecond),
				rejectAfter(errBoom, 100*time.Millisecond),
			},
			want: []promiseallsettled.Result[int]{
				{Val: 30},
				{Err: errBoom},
			},
			elapsed: 200 * time.Millisecond,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()
				got := promiseAllSettled(tt.functions)
				if elapsed := time.Since(start); elapsed > tt.elapsed+threshold {
					t.Errorf("PromiseAllSettled() took %v, want %v", elapsed, tt.elapsed)
				}
				if !slices.Equal(got, tt.want) {
					t.Errorf("PromiseAllSettled() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

func Benchmark(b *testing.B, promiseAllSettled PromiseAllSettledFunc) {
	b.Helper()
	b.ReportAllocs()

	functions := []func() promise.Promiser[int]{
		func() promise.Promiser[int] { return promise.Resolve(1) },
		func() promise.Promiser[int] { return promise.Resolve(2) },
		func() promise.Promiser[int] { return promise.Resolve(3) },
	}

	for b.Loop() {
		if got := promiseAllSettled(functions); len(got) != 3 {
			b.Fatalf("PromiseAllSettled() len = %d", len(got))
		}
	}
}
