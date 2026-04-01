//go:build challenge

package promise_all_settled

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func PromiseAllSettled[T any](functions []func() promise.Promiser[T]) []Result[T] {
	panic("unimplemented")
}
