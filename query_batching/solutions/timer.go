// Package solutions contains implementations of the Query Batching problem.
package solutions

import (
	"sync"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type QueryBatcher struct {
	queryFn func([]string) promise.Promiser[[]string]
	t       time.Duration

	mu       sync.Mutex
	lastTime time.Time
	timer    *time.Timer
	pending  []pending
}

type pending struct {
	key     string
	resolve func(string)
	reject  func(error)
}

func NewQueryBatcher(queryFn func([]string) promise.Promiser[[]string], t time.Duration) *QueryBatcher {
	return &QueryBatcher{
		queryFn: queryFn,
		t:       t,
	}
}

func (qb *QueryBatcher) GetValue(key string) promise.Promiser[string] {
	return promise.New(func(resolve func(string), reject func(error)) {
		qb.mu.Lock()
		now := time.Now()
		if qb.lastTime.IsZero() || !now.Before(qb.lastTime.Add(qb.t)) {
			qb.lastTime = now
			qb.mu.Unlock()

			go func() {
				results, err := qb.queryFn([]string{key}).Result()
				if err != nil {
					reject(err)
					return
				}
				resolve(results[0])
			}()
			return
		}

		qb.pending = append(qb.pending, pending{key: key, resolve: resolve, reject: reject})
		if qb.timer == nil {
			wait := qb.lastTime.Add(qb.t).Sub(now)
			qb.timer = time.AfterFunc(wait, qb.flush)
		}
		qb.mu.Unlock()
	})
}

func (qb *QueryBatcher) flush() {
	qb.mu.Lock()
	items := append([]pending(nil), qb.pending...)
	qb.pending = nil
	qb.timer = nil
	qb.lastTime = time.Now()
	qb.mu.Unlock()

	keys := make([]string, len(items))
	for i, item := range items {
		keys[i] = item.key
	}

	results, err := qb.queryFn(keys).Result()
	if err != nil {
		for _, item := range items {
			item.reject(err)
		}
		return
	}

	for i, item := range items {
		item.resolve(results[i])
	}
}
