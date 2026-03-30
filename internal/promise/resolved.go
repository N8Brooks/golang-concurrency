package promise

import "context"

type Resolved[T any] struct {
	val T
}

func (p *Resolved[T]) Result() (T, error) {
	return p.val, nil
}

func (p *Resolved[T]) ResultContext(_ context.Context) (T, error) {
	return p.val, nil
}

func (p *Resolved[T]) Await() T {
	return p.val
}

// Resolve creates a new Promise that is already resolved with val.
func Resolve[T any](val T) *Resolved[T] {
	return &Resolved[T]{val: val}
}
