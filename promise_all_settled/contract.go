// Package promise_all_settled contains the Promise.allSettled style problem
// and shared result type.
package promise_all_settled

type Obj struct {
	Status string
	Value  int
	Reason string
}
