// Package solutions contains implementations of the multiplex problem.
package solutions

import "context"

type BufferedChannel struct {
	active chan struct{}
}

func NewBufferedChannel(limit int) *BufferedChannel {
	return &BufferedChannel{
		active: make(chan struct{}, limit),
	}
}

func (m *BufferedChannel) Run(ctx context.Context, criticalSection func()) {
	select {
	case m.active <- struct{}{}:
	case <-ctx.Done():
		return
	}
	defer func() {
		<-m.active
	}()

	criticalSection()
}
