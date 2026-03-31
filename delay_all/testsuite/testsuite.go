// Package testsuite contains reusable behavioral tests for delay-all implementations.
package testsuite

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type DelayAllFunc[T comparable] func([]func() promise.Promiser[T], time.Duration) []func() promise.Promiser[T]

func Run(t *testing.T, delayAll DelayAllFunc[int]) {
	t.Helper()

	t.Run("DelayResolution", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			start := time.Now()
			functions := []func() promise.Promiser[int]{
				func() promise.Promiser[int] {
					return promise.New(func(resolve func(int), reject func(error)) {
						time.Sleep(30 * time.Millisecond)
						resolve(10)
					})
				},
			}

			delayed := delayAll(functions, 50*time.Millisecond)
			result := delayed[0]()
			synctest.Wait()
			if got, err := result.Result(); err != nil || got != 10 {
				t.Fatalf("DelayAll(resolve) = %d, %v, want 10, nil", got, err)
			}
			if elapsed := time.Since(start); elapsed != 80*time.Millisecond {
				t.Fatalf("DelayAll(resolve) took %v, want 80ms", elapsed)
			}
		})
	})

	t.Run("DelayRejection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			start := time.Now()
			errBoom := errors.New("boom")
			functions := []func() promise.Promiser[int]{
				func() promise.Promiser[int] {
					return promise.New(func(resolve func(int), reject func(error)) {
						time.Sleep(20 * time.Millisecond)
						reject(errBoom)
					})
				},
			}

			delayed := delayAll(functions, 40*time.Millisecond)
			result := delayed[0]()
			synctest.Wait()
			if _, err := result.Result(); !errors.Is(err, errBoom) {
				t.Fatalf("DelayAll(reject) = %v, want %v", err, errBoom)
			}
			if elapsed := time.Since(start); elapsed != 60*time.Millisecond {
				t.Fatalf("DelayAll(reject) took %v, want 60ms", elapsed)
			}
		})
	})
}
