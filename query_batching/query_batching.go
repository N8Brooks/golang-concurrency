//go:build challenge

// Package query_batching contains the challenge version of the Query Batching
// problem.
package query_batching

import (
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type QueryBatcher struct{}

func NewQueryBatcher(queryFn func([]string) promise.Promiser[[]string], t time.Duration) *QueryBatcher {
	return &QueryBatcher{}
}

func (qb *QueryBatcher) GetValue(key string) promise.Promiser[string] {
	panic("unimplemented")
}
