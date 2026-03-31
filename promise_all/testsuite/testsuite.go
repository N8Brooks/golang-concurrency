// Package testsuite contains reusable behavioral tests for promise_all implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type PromiseAllFunc[T comparable] func(context.Context, []func(context.Context) promise.Promiser[T]) promise.Promiser[[]T]

type testCase[T comparable] struct {
	functions []func(context.Context) promise.Promiser[T]
	res       []T
	err       error
}

func (tc testCase[T]) run(t *testing.T, ctx context.Context, promiseAll PromiseAllFunc[T]) {
	t.Helper()
	p := promiseAll(ctx, tc.functions)
	synctest.Wait()
	if res, err := p.Result(); !slices.Equal(res, tc.res) || !errors.Is(err, tc.err) {
		t.Errorf("PromiseAll(...) = %v, %v, want %v, %v", res, err, tc.res, tc.err)
	}
}

type runnable interface {
	run(t *testing.T, ctx context.Context, promiseAll PromiseAllFunc[int])
}

func resolveAfter[T comparable](res T, d time.Duration) func(context.Context) promise.Promiser[T] {
	return func(ctx context.Context) promise.Promiser[T] {
		return promise.New(func(resolve func(T), reject func(error)) {
			go func() {
				select {
				case <-ctx.Done():
					reject(ctx.Err())
				case <-time.After(d):
					resolve(res)
				}
			}()
		})
	}
}

func rejectAfter[T comparable](err error, d time.Duration) func(context.Context) promise.Promiser[T] {
	return func(ctx context.Context) promise.Promiser[T] {
		return promise.New(func(resolve func(T), reject func(error)) {
			go func() {
				select {
				case <-ctx.Done():
					reject(ctx.Err())
				case <-time.After(d):
					reject(err)
				}
			}()
		})
	}
}

func Run(t *testing.T, promiseAll PromiseAllFunc[int]) {
	t.Helper()

	errBoom := errors.New("error")
	testCases := []runnable{
		testCase[int]{
			functions: []func(context.Context) promise.Promiser[int]{
				resolveAfter(5, 100*time.Millisecond),
			},
			res: []int{5},
		},
		testCase[int]{
			functions: []func(context.Context) promise.Promiser[int]{
				resolveAfter(5, 200*time.Millisecond),
				rejectAfter[int](errBoom, 100*time.Millisecond),
			},
			err: errBoom,
		},
		testCase[int]{
			functions: []func(context.Context) promise.Promiser[int]{
				resolveAfter(4, 50*time.Millisecond),
				resolveAfter(10, 150*time.Millisecond),
				resolveAfter(16, 100*time.Millisecond),
			},
			res: []int{4, 10, 16},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				tc.run(t, t.Context(), promiseAll)
			})
		})
	}

	t.Run("Cancel", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()

			functions := []func(context.Context) promise.Promiser[int]{
				resolveAfter(5, 100*time.Millisecond),
			}

			p := promiseAll(ctx, functions)
			synctest.Wait()
			if _, err := p.Result(); !errors.Is(err, context.Canceled) {
				t.Fatalf("PromiseAll(cancelled) = %v, want %v", err, context.Canceled)
			}
		})
	})
}

func Benchmark[T comparable](b *testing.B, promiseAll PromiseAllFunc[T]) {
	b.Helper()
	b.ReportAllocs()

	functions := []func(context.Context) promise.Promiser[T]{
		func(context.Context) promise.Promiser[T] { return promise.Resolve(*new(T)) },
		func(context.Context) promise.Promiser[T] { return promise.Resolve(*new(T)) },
		func(context.Context) promise.Promiser[T] { return promise.Resolve(*new(T)) },
	}

	for b.Loop() {
		if _, err := promiseAll(b.Context(), functions).Result(); err != nil {
			b.Fatalf("PromiseAll() returned %v", err)
		}
	}
}
