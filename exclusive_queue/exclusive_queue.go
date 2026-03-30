//go:build challenge

// Package exclusive_queue contains the challenge version of the exclusive
// leader/follower queue problem.
//
// Leaders and followers arrive independently. A leader should pair with a
// waiting follower if one exists; otherwise it waits. Followers behave
// symmetrically. Unlike the basic queue problem, each matched leader and
// follower must dance exclusively with each other before any subsequent pair
// may begin dancing.
package exclusive_queue

import "context"

type ExclusiveQueue struct{}

func NewExclusiveQueue() *ExclusiveQueue {
	return &ExclusiveQueue{}
}

func (q *ExclusiveQueue) Leader(ctx context.Context, leader, dance func()) error {
	panic("unimplemented")
}

func (q *ExclusiveQueue) Follower(ctx context.Context, follower, dance func()) error {
	panic("unimplemented")
}
