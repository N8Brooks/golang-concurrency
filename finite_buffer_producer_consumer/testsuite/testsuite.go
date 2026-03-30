// Package testsuite contains reusable behavioral tests for finite-buffer producer-consumer implementations.
package testsuite

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

const (
	defaultCapacity = 2

	eventCreatedBit uint32 = 1 << iota
	eventProcessedBit
)

type ProducerConsumer interface {
	Produce(ctx context.Context, waitForEvent func() int) error
	Consume(ctx context.Context, process func(int)) error
}

func runOneShot(t *testing.T, pc ProducerConsumer, ctx context.Context, consumerFirst bool) {
	t.Helper()

	var state atomic.Uint32
	const expectedEvent = 7

	producerDone := make(chan error, 1)
	consumerDone := make(chan error, 1)

	runProducer := func() {
		producerDone <- pc.Produce(ctx, func() int {
			if state.Load()&eventProcessedBit != 0 {
				t.Error("event creation ran after processing completed")
			}
			state.Or(eventCreatedBit)
			return expectedEvent
		})
	}

	runConsumer := func() {
		consumerDone <- pc.Consume(ctx, func(event int) {
			if event != expectedEvent {
				t.Errorf("consumer got %d, want %d", event, expectedEvent)
			}
			if state.Or(eventProcessedBit)&eventCreatedBit == 0 {
				t.Error("processing ran before the producer created an event")
			}
		})
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
		if err := pc.Produce(ctx, func() int {
			if state.Load()&eventProcessedBit != 0 {
				t.Error("event creation ran after processing completed")
			}
			state.Or(eventCreatedBit)
			return expectedEvent
		}); err != nil {
			t.Fatalf("producer returned %v", err)
		}

		go runConsumer()
	}

	synctest.Wait()

	if consumerFirst {
		select {
		case err := <-producerDone:
			if err != nil {
				t.Fatalf("producer returned %v", err)
			}
		default:
			t.Fatal("producer did not complete")
		}
	}

	select {
	case err := <-consumerDone:
		if err != nil {
			t.Fatalf("consumer returned %v", err)
		}
	default:
		t.Fatal("consumer did not complete")
	}

	if state.Load()&eventProcessedBit == 0 {
		t.Fatal("consumer did not process an event")
	}
}

func Run(t *testing.T, newImpl func(capacity int) ProducerConsumer) {
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
					runOneShot(t, newImpl(defaultCapacity), t.Context(), tc.consumerFirst)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(defaultCapacity)
			ctx := t.Context()

			for i := range 5 {
				runOneShot(t, pc, ctx, i%2 == 0)
			}
		})
	})

	t.Run("MultipleBufferedItems", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(defaultCapacity)
			ctx := t.Context()

			for _, event := range []int{1, 2} {
				if err := pc.Produce(ctx, func() int { return event }); err != nil {
					t.Fatalf("producer returned %v", err)
				}
			}

			var got []int
			for range 2 {
				if err := pc.Consume(ctx, func(event int) {
					got = append(got, event)
				}); err != nil {
					t.Fatalf("consumer returned %v", err)
				}
			}

			slices.Sort(got)
			if !slices.Equal(got, []int{1, 2}) {
				t.Fatalf("consumers got %v, want [1 2] in some order", got)
			}
		})
	})

	t.Run("BufferFullBlocksProducer", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(1)
			ctx := t.Context()

			if err := pc.Produce(ctx, func() int { return 1 }); err != nil {
				t.Fatalf("first producer returned %v", err)
			}

			secondEntered := make(chan struct{})
			secondDone := make(chan error, 1)

			go func() {
				secondDone <- pc.Produce(ctx, func() int {
					close(secondEntered)
					return 2
				})
			}()

			synctest.Wait()

			select {
			case <-secondEntered:
			default:
				t.Fatal("second producer did not start waitForEvent while the buffer was full")
			}

			select {
			case err := <-secondDone:
				t.Fatalf("second producer completed despite the full buffer: %v", err)
			default:
			}

			if err := pc.Consume(ctx, func(event int) {
				if event != 1 {
					t.Errorf("consumer got %d, want 1", event)
				}
			}); err != nil {
				t.Fatalf("consumer returned %v", err)
			}

			synctest.Wait()

			select {
			case err := <-secondDone:
				if err != nil {
					t.Fatalf("second producer returned %v", err)
				}
			default:
				t.Fatal("second producer did not complete after space became available")
			}

			if err := pc.Consume(ctx, func(event int) {
				if event != 2 {
					t.Errorf("consumer got %d, want 2", event)
				}
			}); err != nil {
				t.Fatalf("consumer returned %v", err)
			}
		})
	})

	t.Run("WaitForEventBeforeWaitingForSpace", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(1)
			ctx := t.Context()

			if err := pc.Produce(ctx, func() int { return 1 }); err != nil {
				t.Fatalf("initial producer returned %v", err)
			}

			entered := make(chan struct{})
			release := make(chan struct{})
			done := make(chan error, 1)

			go func() {
				done <- pc.Produce(ctx, func() int {
					close(entered)
					<-release
					return 2
				})
			}()

			synctest.Wait()

			select {
			case <-entered:
			default:
				t.Fatal("producer did not begin waitForEvent before waiting for space")
			}

			close(release)
			synctest.Wait()

			select {
			case err := <-done:
				t.Fatalf("producer completed before a consumer freed space: %v", err)
			default:
			}

			if err := pc.Consume(ctx, func(int) {}); err != nil {
				t.Fatalf("consumer returned %v", err)
			}

			synctest.Wait()

			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("producer returned %v", err)
				}
			default:
				t.Fatal("producer did not complete after space became available")
			}
		})
	})

	t.Run("WaitForEventOutsideCriticalSection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(defaultCapacity)
			ctx := t.Context()

			firstEntered := make(chan struct{})
			firstRelease := make(chan struct{})
			firstDone := make(chan error, 1)
			secondDone := make(chan error, 1)

			go func() {
				firstDone <- pc.Produce(ctx, func() int {
					close(firstEntered)
					<-firstRelease
					return 1
				})
			}()

			synctest.Wait()

			select {
			case <-firstEntered:
			default:
				t.Fatal("first producer did not start waitForEvent")
			}

			go func() {
				secondDone <- pc.Produce(ctx, func() int { return 2 })
			}()

			synctest.Wait()

			select {
			case err := <-secondDone:
				if err != nil {
					t.Fatalf("second producer returned %v", err)
				}
			default:
				t.Fatal("second producer blocked while another producer was still in waitForEvent")
			}

			close(firstRelease)
			synctest.Wait()

			select {
			case err := <-firstDone:
				if err != nil {
					t.Fatalf("first producer returned %v", err)
				}
			default:
				t.Fatal("first producer did not complete")
			}
		})
	})

	t.Run("ProcessOutsideCriticalSection", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(1)
			ctx := t.Context()

			if err := pc.Produce(ctx, func() int { return 1 }); err != nil {
				t.Fatalf("producer returned %v", err)
			}

			processEntered := make(chan struct{})
			processRelease := make(chan struct{})
			firstDone := make(chan error, 1)
			secondDone := make(chan error, 1)

			go func() {
				firstDone <- pc.Consume(ctx, func(event int) {
					if event != 1 {
						t.Errorf("first consumer got %d, want 1", event)
					}
					close(processEntered)
					<-processRelease
				})
			}()

			synctest.Wait()

			select {
			case <-processEntered:
			default:
				t.Fatal("consumer did not begin processing the first item")
			}

			go func() {
				secondDone <- pc.Produce(ctx, func() int { return 2 })
			}()

			synctest.Wait()

			select {
			case err := <-secondDone:
				if err != nil {
					t.Fatalf("second producer returned %v", err)
				}
			default:
				t.Fatal("producer blocked while another consumer was processing outside the buffer")
			}

			close(processRelease)
			synctest.Wait()

			select {
			case err := <-firstDone:
				if err != nil {
					t.Fatalf("first consumer returned %v", err)
				}
			default:
				t.Fatal("first consumer did not complete")
			}

			if err := pc.Consume(ctx, func(event int) {
				if event != 2 {
					t.Errorf("second consumer got %d, want 2", event)
				}
			}); err != nil {
				t.Fatalf("second consumer returned %v", err)
			}
		})
	})

	t.Run("CancelWhileWaitingConsumer", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			pc := newImpl(defaultCapacity)
			var processed atomic.Bool
			done := make(chan error, 1)

			go func() {
				done <- pc.Consume(ctx, func(int) {
					processed.Store(true)
				})
			}()

			synctest.Wait()

			cancel()
			synctest.Wait()

			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("consumer returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("consumer did not exit after cancellation")
			}

			if processed.Load() {
				t.Fatal("consumer processed an item despite cancellation")
			}
		})
	})

	t.Run("CancelWhileWaitingProducer", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(1)

			if err := pc.Produce(t.Context(), func() int { return 1 }); err != nil {
				t.Fatalf("initial producer returned %v", err)
			}

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			var created atomic.Bool
			done := make(chan error, 1)

			go func() {
				done <- pc.Produce(ctx, func() int {
					created.Store(true)
					return 2
				})
			}()

			synctest.Wait()

			if !created.Load() {
				t.Fatal("producer did not begin waitForEvent while waiting for space")
			}

			select {
			case err := <-done:
				t.Fatalf("producer completed despite the full buffer: %v", err)
			default:
			}

			cancel()
			synctest.Wait()

			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("producer returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("producer did not exit after cancellation")
			}

			var got int
			if err := pc.Consume(t.Context(), func(event int) {
				got = event
			}); err != nil {
				t.Fatalf("consumer returned %v", err)
			}
			if got != 1 {
				t.Fatalf("consumer got %d, want 1", got)
			}

			if err := pc.Produce(t.Context(), func() int { return 3 }); err != nil {
				t.Fatalf("replacement producer returned %v", err)
			}
			if err := pc.Consume(t.Context(), func(event int) {
				got = event
			}); err != nil {
				t.Fatalf("replacement consumer returned %v", err)
			}
			if got != 3 {
				t.Fatalf("replacement consumer got %d, want 3", got)
			}
		})
	})

	t.Run("CancelRemovesConsumerWaiter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			pc := newImpl(defaultCapacity)

			waitingCtx, cancel := context.WithCancel(t.Context())
			defer cancel()

			waitingDone := make(chan error, 1)
			go func() {
				waitingDone <- pc.Consume(waitingCtx, func(int) {
					t.Error("canceled consumer processed an item")
				})
			}()

			synctest.Wait()
			cancel()
			synctest.Wait()

			select {
			case err := <-waitingDone:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("waiting consumer returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("waiting consumer did not exit after cancellation")
			}

			if err := pc.Produce(t.Context(), func() int { return 11 }); err != nil {
				t.Fatalf("producer returned %v", err)
			}

			var got int
			if err := pc.Consume(t.Context(), func(event int) {
				got = event
			}); err != nil {
				t.Fatalf("replacement consumer returned %v", err)
			}

			if got != 11 {
				t.Fatalf("replacement consumer got %d, want 11", got)
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func(capacity int) ProducerConsumer) {
	b.Helper()
	b.ReportAllocs()

	b.Run("ProducerFirst", func(b *testing.B) {
		ctx := b.Context()
		pc := newImpl(1)

		for b.Loop() {
			if err := pc.Produce(ctx, func() int { return 1 }); err != nil {
				b.Fatalf("producer returned %v", err)
			}
			if err := pc.Consume(ctx, func(int) {}); err != nil {
				b.Fatalf("consumer returned %v", err)
			}
		}
	})

	b.Run("ConsumerFirst", func(b *testing.B) {
		ctx := b.Context()
		pc := newImpl(1)

		for b.Loop() {
			done := make(chan error, 1)
			go func() {
				done <- pc.Consume(ctx, func(int) {})
			}()

			if err := pc.Produce(ctx, func() int { return 1 }); err != nil {
				b.Fatalf("producer returned %v", err)
			}
			if err := <-done; err != nil {
				b.Fatalf("consumer returned %v", err)
			}
		}
	})
}
