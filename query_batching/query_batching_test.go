//go:build challenge

package query_batching_test

import (
	"testing"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	querybatching "github.com/N8Brooks/golang-concurrency/query_batching"
	"github.com/N8Brooks/golang-concurrency/query_batching/testsuite"
)

func TestQueryBatcher(t *testing.T) {
	testsuite.Run(t, func(queryFn func([]string) promise.Promiser[[]string], d time.Duration) testsuite.Batcher {
		return querybatching.NewQueryBatcher(queryFn, d)
	})
}
