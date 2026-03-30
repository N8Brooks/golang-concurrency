// Package solutions contains implementations of the bounded blocking queue problem.
package solutions

type Channel struct {
	queue chan int
}

func NewChannel(capacity int) *Channel {
	return &Channel{
		queue: make(chan int, capacity),
	}
}

func (q *Channel) Enqueue(element int) {
	q.queue <- element
}

func (q *Channel) Dequeue() int {
	return <-q.queue
}

func (q *Channel) Size() int {
	return len(q.queue)
}
