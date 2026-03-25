package boundedblockingqueue

import (
	"slices"
	"sync"
	"testing"
	"time"
)

func TestBoundedBlockingQueueExample1(t *testing.T) {
	q := Constructor(2)

	done := make(chan int, 2)
	go func() {
		done <- q.Dequeue()
	}()

	q.Enqueue(1)
	first := <-done

	go func() {
		done <- q.Dequeue()
	}()

	q.Enqueue(0)
	second := <-done

	q.Enqueue(2)
	q.Enqueue(3)

	enqueued := make(chan struct{})
	go func() {
		q.Enqueue(4)
		close(enqueued)
	}()

	select {
	case <-enqueued:
		t.Fatal("enqueue should block while queue is full")
	case <-time.After(20 * time.Millisecond):
	}

	third := q.Dequeue()

	select {
	case <-enqueued:
	case <-time.After(time.Second):
		t.Fatal("enqueue should unblock after dequeue")
	}

	if actual := []int{first, second, third, q.Size()}; !slices.Equal(actual, []int{1, 0, 2, 2}) {
		t.Fatalf("got %v, want [1 0 2 2]", actual)
	}
}

func TestBoundedBlockingQueueExample2(t *testing.T) {
	q := Constructor(3)

	start := make(chan struct{})
	var producers sync.WaitGroup
	producers.Add(4)
	for _, element := range []int{1, 0, 2, 3} {
		go func(element int) {
			defer producers.Done()
			<-start
			q.Enqueue(element)
		}(element)
	}

	results := make(chan int, 3)
	var consumers sync.WaitGroup
	consumers.Add(3)
	for range 3 {
		go func() {
			defer consumers.Done()
			<-start
			results <- q.Dequeue()
		}()
	}

	close(start)
	consumers.Wait()

	actual := []int{<-results, <-results, <-results}
	producers.Wait()
	if size := q.Size(); size != 1 {
		t.Fatalf("size = %d, want 1", size)
	}

	actual = append(actual, q.Dequeue())
	slices.Sort(actual)
	if !slices.Equal(actual, []int{0, 1, 2, 3}) {
		t.Fatalf("got %v, want a permutation of [0 1 2 3]", actual)
	}
}

func TestBoundedBlockingQueueSize(t *testing.T) {
	q := Constructor(2)
	if size := q.Size(); size != 0 {
		t.Fatalf("size = %d, want 0", size)
	}

	q.Enqueue(7)
	q.Enqueue(8)
	if size := q.Size(); size != 2 {
		t.Fatalf("size = %d, want 2", size)
	}

	if value := q.Dequeue(); value != 7 {
		t.Fatalf("dequeue = %d, want 7", value)
	}
	if size := q.Size(); size != 1 {
		t.Fatalf("size = %d, want 1", size)
	}
}
