package boundedblockingqueue

type BoundedBlockingQueue struct {
	queue chan int
}

func Constructor(capacity int) *BoundedBlockingQueue {
	return &BoundedBlockingQueue{
		queue: make(chan int, capacity),
	}
}

func (q *BoundedBlockingQueue) Enqueue(element int) {
	q.queue <- element
}

func (q *BoundedBlockingQueue) Dequeue() int {
	return <-q.queue
}

func (q *BoundedBlockingQueue) Size() int {
	return len(q.queue)
}
