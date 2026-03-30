package solutions

import (
	"context"
	"sync"
)

type consumerWaiter struct {
	event chan int
}

type Channel struct {
	mu        sync.Mutex
	buffer    []int
	consumers []*consumerWaiter
}

func NewChannel() *Channel {
	return &Channel{}
}

func (pc *Channel) Produce(_ context.Context, waitForEvent func() int) error {
	event := waitForEvent()

	pc.mu.Lock()
	if len(pc.consumers) > 0 {
		waiter := pc.consumers[0]
		copy(pc.consumers, pc.consumers[1:])
		pc.consumers[len(pc.consumers)-1] = nil
		pc.consumers = pc.consumers[:len(pc.consumers)-1]
		pc.mu.Unlock()

		waiter.event <- event
		return nil
	}
	pc.buffer = append(pc.buffer, event)
	pc.mu.Unlock()
	return nil
}

func (pc *Channel) Consume(ctx context.Context, process func(int)) error {
	pc.mu.Lock()
	if len(pc.buffer) > 0 {
		event := pc.buffer[0]
		pc.buffer = pc.buffer[1:]
		pc.mu.Unlock()
		process(event)
		return nil
	}
	if err := ctx.Err(); err != nil {
		pc.mu.Unlock()
		return err
	}

	waiter := &consumerWaiter{
		event: make(chan int, 1),
	}
	pc.consumers = append(pc.consumers, waiter)
	pc.mu.Unlock()

	select {
	case event := <-waiter.event:
		process(event)
		return nil
	case <-ctx.Done():
		pc.mu.Lock()
		index := -1
		for i, candidate := range pc.consumers {
			if candidate == waiter {
				index = i
				break
			}
		}
		if index >= 0 {
			copy(pc.consumers[index:], pc.consumers[index+1:])
			pc.consumers[len(pc.consumers)-1] = nil
			pc.consumers = pc.consumers[:len(pc.consumers)-1]
			pc.mu.Unlock()
			return ctx.Err()
		}
		pc.mu.Unlock()

		event := <-waiter.event
		process(event)
		return nil
	}
}
