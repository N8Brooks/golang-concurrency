// Package promise provides a small future/promise helper used by async
// problems in this repository.
//
// The implementation is adapted from the future example in cloud-native-go:
// https://github.com/cloud-native-go/examples/blob/main/ch04/future.go
package promise

import (
	"context"
	"sync"
)

// Promiser is a future type that is settled asynchronously.
type Promiser[T any] interface {
	// Result returns the value exclusive or err of the Promise asynchronously.
	Result() (T, error)
	// ResultContext returns the value exclusive or err of the Promise with respect to the context.
	// If the context is canceled before the Promise is settled, the error is the context error.
	// This will be cached and returned in further calls of Result or ResultContext.
	ResultContext(ctx context.Context) (T, error)
	// Await returns the value of the Promise asynchronously.
	// NOTE: Await panics if the Promise is rejected.
	Await() T
}

// Promise is a promise that is settled asynchronously.
type Promise[T any] struct {
	// Used to ensure the promise is settled only once
	once sync.Once
	// The cached result of the promise
	result[T]
	// Used to pass the result to the Result method asynchronously
	ch <-chan result[T]
}

func (p *Promise[T]) Result() (T, error) {
	p.once.Do(func() {
		p.result = <-p.ch
	})
	return p.val, p.err
}

func (p *Promise[T]) ResultContext(ctx context.Context) (T, error) {
	p.once.Do(func() {
		select {
		case p.result = <-p.ch:
		case <-ctx.Done():
			p.err = ctx.Err()
		}
	})
	return p.val, p.err
}

func (p *Promise[T]) Await() T {
	val, err := p.Result()
	if err != nil {
		panic(err)
	}
	return val
}

// New creates a new Promise with executer which receives callbacks for resolving the promise.
// The Promise is settled with the first call to resolve or reject. Further calls will be ignored.
// NOTE: reject always panics if called with nil. The Promise will still be able to be settled with the callbacks.
func New[T any](executer func(resolve func(T), reject func(error))) *Promise[T] {
	var once sync.Once
	ch := make(chan result[T], 1)

	settle := func(res result[T]) {
		once.Do(func() {
			defer close(ch)
			ch <- res
		})
	}

	resolve := func(val T) {
		settle(result[T]{val: val})
	}

	reject := func(err error) {
		if err == nil {
			panic("reject called with nil")
		}
		settle(result[T]{err: err})
	}

	executer(resolve, reject)

	return &Promise[T]{ch: ch}
}

// result is the de-normalized discriminated union of the promise result
type result[T any] struct {
	val T
	err error
}
