package promise

import "context"

type Rejected[T any] struct {
	err error
}

// Reject creates a new Promise that is already rejected with err.
// NOTE: Reject panics if called with nil.
func Reject[T any](err error) *Rejected[T] {
	if err == nil {
		panic("reject called with nil")
	}
	return &Rejected[T]{err: err}
}

func (p *Rejected[T]) Result() (T, error) {
	var val T
	return val, p.err
}

func (p *Rejected[T]) ResultContext(_ context.Context) (T, error) {
	var val T
	return val, p.err
}

func (p *Rejected[T]) Await() T {
	panic(p.err)
}
