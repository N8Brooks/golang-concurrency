// Package testsuite contains reusable behavioral tests for no-starve
// readers-writers implementations.
package testsuite

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

type NoStarveReadersWriters interface {
	Reader(ctx context.Context, read func()) error
	Writer(ctx context.Context, write func()) error
}

type worker struct {
	entered chan struct{}
	release chan struct{}
	done    chan error
}

func startReader(rw NoStarveReadersWriters, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- rw.Reader(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()

	return w
}

func startWriter(rw NoStarveReadersWriters, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		w.done <- rw.Writer(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()

	return w
}

func requireClosed(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
	default:
		t.Fatal(message)
	}
}

func requireOpen(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
		t.Fatal(message)
	default:
	}
}

func requireSuccess(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("%s: %v", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func requireCanceled(t *testing.T, ch <-chan error, message string) {
	t.Helper()

	select {
	case err := <-ch:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: got %v, want context.Canceled", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func runReadersShare(t *testing.T, rw NoStarveReadersWriters, ctx context.Context) {
	t.Helper()

	r1 := startReader(rw, ctx)
	synctest.Wait()
	requireClosed(t, r1.entered, "first reader did not enter")

	r2 := startReader(rw, ctx)
	synctest.Wait()
	requireClosed(t, r2.entered, "second reader did not enter alongside first reader")

	close(r1.release)
	close(r2.release)
	synctest.Wait()

	requireSuccess(t, r1.done, "first reader did not complete")
	requireSuccess(t, r2.done, "second reader did not complete")
}

func runWriterExclusive(t *testing.T, rw NoStarveReadersWriters, ctx context.Context) {
	t.Helper()

	w1 := startWriter(rw, ctx)
	synctest.Wait()
	requireClosed(t, w1.entered, "first writer did not enter")

	w2 := startWriter(rw, ctx)
	synctest.Wait()
	requireOpen(t, w2.entered, "second writer entered while first writer was active")

	close(w1.release)
	synctest.Wait()

	requireSuccess(t, w1.done, "first writer did not complete")
	requireClosed(t, w2.entered, "second writer did not enter after first writer left")

	close(w2.release)
	synctest.Wait()

	requireSuccess(t, w2.done, "second writer did not complete")
}

func Run(t *testing.T, newImpl func() NoStarveReadersWriters) {
	t.Helper()

	t.Run("ReadersShare", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runReadersShare(t, newImpl(), t.Context())
		})
	})

	t.Run("WriterWaitsForReaders", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			ctx := t.Context()

			r1 := startReader(rw, ctx)
			r2 := startReader(rw, ctx)
			synctest.Wait()
			requireClosed(t, r1.entered, "first reader did not enter")
			requireClosed(t, r2.entered, "second reader did not enter")

			w := startWriter(rw, ctx)
			synctest.Wait()
			requireOpen(t, w.entered, "writer entered while readers were active")

			close(r1.release)
			close(r2.release)
			synctest.Wait()

			requireSuccess(t, r1.done, "first reader did not complete")
			requireSuccess(t, r2.done, "second reader did not complete")
			requireClosed(t, w.entered, "writer did not enter after readers left")

			close(w.release)
			synctest.Wait()
			requireSuccess(t, w.done, "writer did not complete")
		})
	})

	t.Run("ReadersWaitForWriter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			ctx := t.Context()

			w := startWriter(rw, ctx)
			synctest.Wait()
			requireClosed(t, w.entered, "writer did not enter")

			r1 := startReader(rw, ctx)
			r2 := startReader(rw, ctx)
			synctest.Wait()
			requireOpen(t, r1.entered, "first reader entered while writer was active")
			requireOpen(t, r2.entered, "second reader entered while writer was active")

			close(w.release)
			synctest.Wait()

			requireSuccess(t, w.done, "writer did not complete")
			requireClosed(t, r1.entered, "first reader did not enter after writer left")
			requireClosed(t, r2.entered, "second reader did not enter after writer left")

			close(r1.release)
			close(r2.release)
			synctest.Wait()

			requireSuccess(t, r1.done, "first reader did not complete")
			requireSuccess(t, r2.done, "second reader did not complete")
		})
	})

	t.Run("QueuedWriterBlocksLaterReaders", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			ctx := t.Context()

			activeReader := startReader(rw, ctx)
			synctest.Wait()
			requireClosed(t, activeReader.entered, "active reader did not enter")

			waitingWriter := startWriter(rw, ctx)
			synctest.Wait()
			requireOpen(t, waitingWriter.entered, "writer entered while reader was active")

			lateReader := startReader(rw, ctx)
			synctest.Wait()
			requireOpen(t, lateReader.entered, "late reader bypassed queued writer")

			close(activeReader.release)
			synctest.Wait()

			requireSuccess(t, activeReader.done, "active reader did not complete")
			requireClosed(t, waitingWriter.entered, "queued writer did not enter after readers left")
			requireOpen(t, lateReader.entered, "late reader entered before queued writer finished")

			close(waitingWriter.release)
			synctest.Wait()

			requireSuccess(t, waitingWriter.done, "queued writer did not complete")
			requireClosed(t, lateReader.entered, "late reader did not enter after queued writer finished")

			close(lateReader.release)
			synctest.Wait()
			requireSuccess(t, lateReader.done, "late reader did not complete")
		})
	})

	t.Run("WriterExclusive", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runWriterExclusive(t, newImpl(), t.Context())
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			ctx := t.Context()
			runReadersShare(t, rw, ctx)
			runWriterExclusive(t, rw, ctx)
			runReadersShare(t, rw, ctx)
		})
	})

	t.Run("CancelWhileWaitingReader", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			active := startWriter(rw, t.Context())
			synctest.Wait()
			requireClosed(t, active.entered, "writer did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			reader := startReader(rw, ctx)
			synctest.Wait()
			requireOpen(t, reader.entered, "reader entered while writer was active")

			cancel()
			synctest.Wait()
			requireCanceled(t, reader.done, "reader did not exit after cancellation")

			close(active.release)
			synctest.Wait()
			requireSuccess(t, active.done, "writer did not complete")
		})
	})

	t.Run("CancelWhileWaitingWriter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			active := startReader(rw, t.Context())
			synctest.Wait()
			requireClosed(t, active.entered, "reader did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			writer := startWriter(rw, ctx)
			synctest.Wait()
			requireOpen(t, writer.entered, "writer entered while reader was active")

			cancel()
			synctest.Wait()
			requireCanceled(t, writer.done, "writer did not exit after cancellation")

			close(active.release)
			synctest.Wait()
			requireSuccess(t, active.done, "reader did not complete")
		})
	})

	t.Run("CanceledQueuedWriterDoesNotBlockFutureReaders", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rw := newImpl()
			active := startReader(rw, t.Context())
			synctest.Wait()
			requireClosed(t, active.entered, "reader did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			writer := startWriter(rw, ctx)
			synctest.Wait()
			requireOpen(t, writer.entered, "writer entered while reader was active")

			cancel()
			synctest.Wait()
			requireCanceled(t, writer.done, "writer did not exit after cancellation")

			lateReader := startReader(rw, t.Context())
			synctest.Wait()
			requireClosed(t, lateReader.entered, "late reader did not enter after canceled writer was removed")

			close(active.release)
			synctest.Wait()
			requireSuccess(t, active.done, "reader did not complete")

			close(lateReader.release)
			synctest.Wait()
			requireSuccess(t, lateReader.done, "late reader did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() NoStarveReadersWriters) {
	b.Helper()
	b.ReportAllocs()

	b.Run("ReadOnly", func(b *testing.B) {
		rw := newImpl()
		ctx := b.Context()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := rw.Reader(ctx, func() {}); err != nil {
					b.Fatalf("reader returned %v", err)
				}
			}
		})
	})

	b.Run("Mixed", func(b *testing.B) {
		rw := newImpl()
		ctx := b.Context()
		var counter atomic.Uint64

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if counter.Add(1)%8 == 0 {
					if err := rw.Writer(ctx, func() {}); err != nil {
						b.Fatalf("writer returned %v", err)
					}
					continue
				}
				if err := rw.Reader(ctx, func() {}); err != nil {
					b.Fatalf("reader returned %v", err)
				}
			}
		})
	})
}
