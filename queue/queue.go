//go:build challenge

// Package queue contains the challenge version of the leader/follower queue
// problem.
//
// Leaders and followers arrive independently. When a leader arrives, it should
// pair with a waiting follower if one exists; otherwise it waits. Followers
// behave symmetrically. A matched leader and follower may both proceed to
// dance, and unmatched callers should be able to stop waiting if their context
// is canceled.
package queue

import "context"

type Queue struct{}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Leader(ctx context.Context, leader, dance func()) error {
	panic("unimplemented")
}

func (q *Queue) Follower(ctx context.Context, follower, dance func()) error {
	panic("unimplemented")
}
