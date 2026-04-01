package solutions

import (
	"context"
)

type Channel struct {
	event chan int
}

func NewChannel() *Channel {
	return &Channel{
		event: make(chan int),
	}
}

func (pc *Channel) Produce(ctx context.Context, waitForEvent func() int) {
	select {
	case <-ctx.Done():
		return
	case pc.event <- waitForEvent():
	}
}

func (pc *Channel) Consume(ctx context.Context, process func(int)) {
	select {
	case <-ctx.Done():
		return
	case event := <-pc.event:
		process(event)
	}
}
