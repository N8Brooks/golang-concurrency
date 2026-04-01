//go:build challenge

// Package producerconsumer contains the challenge version of the
// producer-consumer problem.
//
// Producers wait for events and append them to a shared buffer. Consumers
// remove buffered events and process them. Buffer access must be exclusive, a
// consumer must wait while the buffer is empty, and waiting consumers should be
// able to stop waiting when their context is canceled.
package producerconsumer

import (
	"context"
)

type ProducerConsumer struct {
	buffer []int
}

func NewProducerConsumer() *ProducerConsumer {
	return &ProducerConsumer{}
}

func (pc *ProducerConsumer) Produce(_ context.Context, waitForEvent func() int) {
	event := waitForEvent()
	pc.buffer = append(pc.buffer, event)
}

func (pc *ProducerConsumer) Consume(ctx context.Context, process func(int)) {
	n := len(pc.buffer) - 1
	event := pc.buffer[n]
	pc.buffer = pc.buffer[:n]
	process(event)
}
