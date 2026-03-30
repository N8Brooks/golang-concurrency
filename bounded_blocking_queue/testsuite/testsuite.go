// Package testsuite contains reusable behavioral tests for bounded blocking queue implementations.
package testsuite

import (
	"slices"
	"sync"
	"testing"
	"time"
)

type Queue interface {
	Enqueue(element int)
	Dequeue() int
	Size() int
}

func Run(t *testing.T, newImpl func(capacity int) Queue) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		q := newImpl(2)

		q.Enqueue(7)
		if size := q.Size(); size != 1 {
			t.Fatalf("size = %d, want 1", size)
		}
		if got := q.Dequeue(); got != 7 {
			t.Fatalf("dequeue = %d, want 7", got)
		}
		if size := q.Size(); size != 0 {
			t.Fatalf("size = %d, want 0", size)
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		q := newImpl(3)

		for _, value := range []int{1, 2, 3, 4, 5} {
			q.Enqueue(value)
			if got := q.Dequeue(); got != value {
				t.Fatalf("dequeue = %d, want %d", got, value)
			}
		}

		if size := q.Size(); size != 0 {
			t.Fatalf("size = %d, want 0", size)
		}
	})

	t.Run("BlocksWhenEmpty", func(t *testing.T) {
		q := newImpl(2)

		done := make(chan int, 1)
		go func() {
			done <- q.Dequeue()
		}()

		select {
		case got := <-done:
			t.Fatalf("dequeue returned %d while queue was empty", got)
		case <-time.After(20 * time.Millisecond):
		}

		q.Enqueue(11)

		select {
		case got := <-done:
			if got != 11 {
				t.Fatalf("dequeue = %d, want 11", got)
			}
		case <-time.After(time.Second):
			t.Fatal("dequeue did not unblock after enqueue")
		}
	})

	t.Run("BlocksWhenFull", func(t *testing.T) {
		q := newImpl(2)
		q.Enqueue(1)
		q.Enqueue(2)

		enqueued := make(chan struct{})
		go func() {
			q.Enqueue(3)
			close(enqueued)
		}()

		select {
		case <-enqueued:
			t.Fatal("enqueue should block while queue is full")
		case <-time.After(20 * time.Millisecond):
		}

		if got := q.Dequeue(); got != 1 {
			t.Fatalf("dequeue = %d, want 1", got)
		}

		select {
		case <-enqueued:
		case <-time.After(time.Second):
			t.Fatal("enqueue did not unblock after dequeue")
		}

		if got := q.Dequeue(); got != 2 {
			t.Fatalf("dequeue = %d, want 2", got)
		}
		if got := q.Dequeue(); got != 3 {
			t.Fatalf("dequeue = %d, want 3", got)
		}
	})

	t.Run("ConcurrentPermutation", func(t *testing.T) {
		q := newImpl(3)

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
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) Queue) {
	b.Helper()
	b.ReportAllocs()

	b.Run("EnqueueDequeue", func(b *testing.B) {
		q := newImpl(1)
		for i := 0; b.Loop(); i++ {
			q.Enqueue(i)
			if got := q.Dequeue(); got != i {
				b.Fatalf("dequeue = %d, want %d", got, i)
			}
		}
	})
}
