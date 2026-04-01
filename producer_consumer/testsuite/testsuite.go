// Package testsuite contains reusable behavioral tests for producer-consumer implementations.
package testsuite

import (
	"context"
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const (
	eventCreatedBit uint32 = 1 << iota
	eventProcessedBit
)

type ProducerConsumer interface {
	Produce(ctx context.Context, waitForEvent func() int)
	Consume(ctx context.Context, process func(int))
}

func runOneShot(t *testing.T, pc ProducerConsumer, ctx context.Context, consumerFirst bool) {
	t.Helper()

	var state atomic.Uint32
	const expectedEvent = 7

	producerDone := make(chan struct{}, 1)
	consumerDone := make(chan struct{}, 1)

	runProducer := func() {
		pc.Produce(ctx, func() int {
			if state.Or(eventCreatedBit)&eventProcessedBit != 0 {
				t.Error("event creation ran after processing completed")
			}
			return expectedEvent
		})
		producerDone <- struct{}{}
	}

	runConsumer := func() {
		pc.Consume(ctx, func(event int) {
			if event != expectedEvent {
				t.Errorf("consumer got %d, want %d", event, expectedEvent)
			}
			if state.Or(eventProcessedBit)&eventCreatedBit == 0 {
				t.Error("processing ran before the producer created an event")
			}
		})
		consumerDone <- struct{}{}
	}

	if consumerFirst {
		go runConsumer()
		synctest.Wait()

		select {
		case err := <-consumerDone:
			t.Fatalf("consumer completed before an item was produced: %v", err)
		default:
		}

		go runProducer()
	} else {
		go pc.Produce(ctx, func() int {
			if state.Or(eventCreatedBit)&eventProcessedBit != 0 {
				t.Error("event creation ran after processing completed")
			}
			return expectedEvent
		})

		go runConsumer()
	}

	synctest.Wait()

	if consumerFirst {
		select {
		case <-producerDone:
		default:
			t.Fatal("producer did not complete")
		}
	}

	select {
	case <-consumerDone:
	default:
		t.Fatal("consumer did not complete")
	}

	if state.Load()&eventProcessedBit == 0 {
		t.Fatal("consumer did not process an event")
	}
}

func Run(t *testing.T, newImpl func() ProducerConsumer) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, tc := range []struct {
			name          string
			consumerFirst bool
		}{
			{name: "ConsumerFirst", consumerFirst: true},
			{name: "ProducerFirst", consumerFirst: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					runOneShot(t, newImpl(), t.Context(), tc.consumerFirst)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl()
			ctx := t.Context()

			for i := range 5 {
				runOneShot(t, pc, ctx, i%2 == 0)
			}
		})
	})

	t.Run("MultipleBufferedItems", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			expected := []int{1, 2}
			pc := newImpl()
			ctx := t.Context()

			for _, event := range expected {
				go pc.Produce(ctx, func() int { return event })
			}

			actual := make(chan int, len(expected))
			for range len(expected) {
				go pc.Consume(ctx, func(event int) {
					actual <- event
				})
			}

			synctest.Wait()

			var got []int
			for range expected {
				select {
				case event := <-actual:
					got = append(got, event)
				default:
					break
				}
			}
			slices.Sort(got)

			if !slices.Equal(got, []int{1, 2}) {
				t.Fatalf("consumers got %v, want [1 2] in some order", got)
			}
		})
	})

	t.Run("WaitForEventOutsideCriticalSection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl()
			ctx := t.Context()

			firstEntered := make(chan struct{})
			firstRelease := make(chan struct{})
			firstDone := make(chan struct{})
			secondEntered := make(chan struct{})
			secondDone := make(chan struct{})

			go func() {
				pc.Produce(ctx, func() int {
					close(firstEntered)
					<-firstRelease
					return 1
				})
				close(firstDone)
			}()

			synctest.Wait()

			select {
			case <-firstEntered:
			default:
				t.Fatal("first producer did not start waitForEvent")
			}

			go func() {
				pc.Produce(ctx, func() int {
					close(secondEntered)
					return 2
				})
				close(secondDone)
			}()

			synctest.Wait()

			select {
			case <-secondEntered:
			default:
				t.Fatal("second producer could not enter waitForEvent while first producer was blocked")
			}

			close(firstRelease)
			go pc.Consume(ctx, func(event int) {})
			go pc.Consume(ctx, func(event int) {})
			synctest.Wait()

			select {
			case <-firstDone:
			default:
				t.Fatal("first producer did not complete")
			}

			select {
			case <-secondDone:
			default:
				t.Fatal("second producer did not complete")
			}
		})
	})

	t.Run("ProcessOutsideCriticalSection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl()
			ctx := t.Context()

			go pc.Produce(ctx, func() int { return 1 })
			go pc.Consume(ctx, func(int) {})
			synctest.Wait()

			processEntered := make(chan struct{})
			processRelease := make(chan struct{})
			firstDone := make(chan struct{})
			secondEntered := make(chan struct{})
			secondDone := make(chan struct{})

			go func() {
				pc.Consume(ctx, func(event int) {
					if event != 1 {
						t.Errorf("first consumer got %d, want 1", event)
					}
					close(processEntered)
					<-processRelease
				})
				close(firstDone)
			}()

			go pc.Produce(ctx, func() int { return 1 })
			synctest.Wait()

			select {
			case <-processEntered:
			default:
				t.Fatal("consumer did not begin processing the first item")
			}

			go func() {
				pc.Produce(ctx, func() int {
					close(secondEntered)
					return 2
				})
				close(secondDone)
			}()

			synctest.Wait()

			select {
			case <-secondEntered:
			default:
				t.Fatal("second producer could not enter waitForEvent while another consumer was processing")
			}

			close(processRelease)
			go pc.Consume(ctx, func(int) {})
			synctest.Wait()

			select {
			case <-firstDone:
			default:
				t.Fatal("first consumer did not complete")
			}

			select {
			case <-secondDone:
			default:
				t.Fatal("second producer did not complete")
			}
		})
	})

	t.Run("CancelWhileWaiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			pc := newImpl()
			var processed atomic.Bool
			done := make(chan struct{}, 1)

			go func() {
				pc.Consume(ctx, func(int) {
					processed.Store(true)
				})
				done <- struct{}{}
			}()

			synctest.Wait()

			cancel()
			synctest.Wait()

			select {
			case <-done:
			default:
				t.Fatal("consumer did not exit after cancellation")
			}

			if processed.Load() {
				t.Fatal("consumer processed an item despite cancellation")
			}
		})
	})

	t.Run("CancelRemovesWaiter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl()

			waitingCtx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waitingDone := make(chan struct{}, 1)
			go func() {
				pc.Consume(waitingCtx, func(int) {
					t.Error("canceled consumer processed an item")
				})
				waitingDone <- struct{}{}
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			select {
			case <-waitingDone:
			default:
				t.Fatal("waiting consumer did not exit after cancellation")
			}

			var got int
			produced := make(chan struct{})
			consumed := make(chan struct{})

			go func() {
				pc.Produce(t.Context(), func() int { return 11 })
				close(produced)
			}()

			go func() {
				pc.Consume(t.Context(), func(event int) {
					got = event
				})
				close(consumed)
			}()

			synctest.Wait()

			select {
			case <-produced:
			default:
				t.Fatal("replacement producer did not complete")
			}
			select {
			case <-consumed:
			default:
				t.Fatal("replacement consumer did not complete")
			}

			if got != 11 {
				t.Fatalf("replacement consumer got %d, want 11", got)
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func() ProducerConsumer) {
	b.Helper()
	b.ReportAllocs()

	b.Run("ProducerFirst", func(b *testing.B) {
		ctx := b.Context()
		pc := newImpl()

		for b.Loop() {
			go pc.Produce(ctx, func() int { return 1 })
			pc.Consume(ctx, func(int) {})
		}
	})

	b.Run("ConsumerFirst", func(b *testing.B) {
		ctx := b.Context()
		pc := newImpl()

		for b.Loop() {
			go pc.Consume(ctx, func(int) {})
			pc.Produce(ctx, func() int { return 1 })
		}
	})
}
