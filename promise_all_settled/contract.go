// Package promise_all_settled contains the Promise.allSettled style problem
// and shared result type.
package promise_all_settled

type Result[T any] struct {
	Val T
	Err error
}
