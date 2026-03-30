package solutions

import "context"

type Channel struct {
	buffer chan int
}

func NewChannel(capacity int) *Channel {
	return &Channel{
		buffer: make(chan int, capacity),
	}
}

func (pc *Channel) Produce(ctx context.Context, waitForEvent func() int) error {
	event := waitForEvent()

	select {
	case pc.buffer <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (pc *Channel) Consume(ctx context.Context, process func(int)) error {
	select {
	case event := <-pc.buffer:
		process(event)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
