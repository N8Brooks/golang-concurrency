// Package testsuite contains reusable behavioral tests for promisify implementations.
package testsuite

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type PromisifyFunc func(fn func(callback func(int, error), args ...int)) func(args ...int) promise.Promiser[int]

func Run(t *testing.T, promisify PromisifyFunc) {
	t.Helper()

	t.Run("Resolve", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fn := func(callback func(int, error), args ...int) {
				callback(args[0]+args[1], nil)
			}

			if got, err := promisify(fn)(2, 5).Result(); err != nil || got != 7 {
				t.Fatalf("Promisify(resolve) = %d, %v, want 7, nil", got, err)
			}
		})
	})

	t.Run("Reject", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			errBoom := errors.New("boom")
			fn := func(callback func(int, error), args ...int) {
				callback(0, errBoom)
			}

			if _, err := promisify(fn)(2, 5).Result(); !errors.Is(err, errBoom) {
				t.Fatalf("Promisify(reject) = %v, want %v", err, errBoom)
			}
		})
	})

	t.Run("AsyncCallback", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fn := func(callback func(int, error), args ...int) {
				go func() {
					time.Sleep(50 * time.Millisecond)
					callback(args[0]*args[1], nil)
				}()
			}

			result := promisify(fn)(3, 4)
			synctest.Wait()
			if got, err := result.Result(); err != nil || got != 12 {
				t.Fatalf("Promisify(async) = %d, %v, want 12, nil", got, err)
			}
		})
	})
}
