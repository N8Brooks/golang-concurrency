// Package solutions contains implementations of the mutex problem.
package solutions

import "context"

type Channel struct {
	token chan struct{}
}

func NewChannel() *Channel {
	token := make(chan struct{}, 1)
	token <- struct{}{}
	return &Channel{token: token}
}

func (m *Channel) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.token:
		return nil
	}
}

func (m *Channel) Unlock() {
	m.token <- struct{}{}
}
