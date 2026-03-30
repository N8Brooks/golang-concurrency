//go:build challenge

package promise_all_settled

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func PromiseAllSettled(functions []func() promise.Promiser[int]) promise.Promiser[[]Obj] {
	panic("unimplemented")
}
