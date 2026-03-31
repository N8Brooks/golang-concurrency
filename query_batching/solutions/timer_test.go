package solutions_test

import (
	"testing"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	"github.com/N8Brooks/golang-concurrency/query_batching/solutions"
	"github.com/N8Brooks/golang-concurrency/query_batching/testsuite"
)

func TestTimer(t *testing.T) {
	testsuite.Run(t, func(queryFn func([]string) promise.Promiser[[]string], d time.Duration) testsuite.Batcher {
		return solutions.NewQueryBatcher(queryFn, d)
	})
}
