//go:build challenge

// Package boundedblockingqueue contains the challenge version of the bounded
// blocking queue problem.
//
// The queue has a fixed capacity. Enqueue blocks while the queue is full,
// Dequeue blocks while it is empty, and Size reports the current number of
// elements.
package boundedblockingqueue

type Queue struct{}

func NewQueue(capacity int) *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(element int) {
	panic("unimplemented")
}

func (q *Queue) Dequeue() int {
	panic("unimplemented")
}

func (q *Queue) Size() int {
	panic("unimplemented")
}
