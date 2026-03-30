// Package testsuite contains reusable behavioral tests for search-insert-delete implementations.
package testsuite

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
)

type SearchInsertDelete interface {
	Searcher(ctx context.Context, search func()) error
	Inserter(ctx context.Context, insert func()) error
	Deleter(ctx context.Context, remove func()) error
}

type worker struct {
	entered chan struct{}
	release chan struct{}
	done    chan error
}

func startSearcher(sid SearchInsertDelete, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}
	go func() {
		w.done <- sid.Searcher(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()
	return w
}

func startInserter(sid SearchInsertDelete, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}
	go func() {
		w.done <- sid.Inserter(ctx, func() {
			close(w.entered)
			<-w.release
		})
	}()
	return w
}

func startDeleter(sid SearchInsertDelete, ctx context.Context) worker {
	w := worker{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}
	go func() {
		w.done <- sid.Deleter(ctx, func() {
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

func runSearchersShare(t *testing.T, sid SearchInsertDelete, ctx context.Context) {
	t.Helper()

	s1 := startSearcher(sid, ctx)
	synctest.Wait()
	requireClosed(t, s1.entered, "first searcher did not enter")

	s2 := startSearcher(sid, ctx)
	synctest.Wait()
	requireClosed(t, s2.entered, "second searcher did not enter alongside first searcher")

	close(s1.release)
	close(s2.release)
	synctest.Wait()

	requireSuccess(t, s1.done, "first searcher did not complete")
	requireSuccess(t, s2.done, "second searcher did not complete")
}

func runInserterExclusive(t *testing.T, sid SearchInsertDelete, ctx context.Context) {
	t.Helper()

	i1 := startInserter(sid, ctx)
	synctest.Wait()
	requireClosed(t, i1.entered, "first inserter did not enter")

	i2 := startInserter(sid, ctx)
	synctest.Wait()
	requireOpen(t, i2.entered, "second inserter entered while first inserter was active")

	close(i1.release)
	synctest.Wait()
	requireSuccess(t, i1.done, "first inserter did not complete")
	requireClosed(t, i2.entered, "second inserter did not enter after first inserter left")

	close(i2.release)
	synctest.Wait()
	requireSuccess(t, i2.done, "second inserter did not complete")
}

func Run(t *testing.T, newImpl func() SearchInsertDelete) {
	t.Helper()

	t.Run("SearchersShare", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runSearchersShare(t, newImpl(), t.Context())
		})
	})

	t.Run("InsertersAreExclusive", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runInserterExclusive(t, newImpl(), t.Context())
		})
	})

	t.Run("SearcherSharesWithInserter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			ctx := t.Context()

			s := startSearcher(sid, ctx)
			synctest.Wait()
			requireClosed(t, s.entered, "searcher did not enter")

			i := startInserter(sid, ctx)
			synctest.Wait()
			requireClosed(t, i.entered, "inserter did not enter alongside searcher")

			close(s.release)
			close(i.release)
			synctest.Wait()

			requireSuccess(t, s.done, "searcher did not complete")
			requireSuccess(t, i.done, "inserter did not complete")
		})
	})

	t.Run("DeleterWaitsForSearchersAndInserters", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			ctx := t.Context()

			s := startSearcher(sid, ctx)
			i := startInserter(sid, ctx)
			synctest.Wait()
			requireClosed(t, s.entered, "searcher did not enter")
			requireClosed(t, i.entered, "inserter did not enter")

			d := startDeleter(sid, ctx)
			synctest.Wait()
			requireOpen(t, d.entered, "deleter entered while searcher and inserter were active")

			close(s.release)
			close(i.release)
			synctest.Wait()

			requireSuccess(t, s.done, "searcher did not complete")
			requireSuccess(t, i.done, "inserter did not complete")
			requireClosed(t, d.entered, "deleter did not enter after searcher and inserter left")

			close(d.release)
			synctest.Wait()
			requireSuccess(t, d.done, "deleter did not complete")
		})
	})

	t.Run("SearcherAndInserterWaitForDeleter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			ctx := t.Context()

			d := startDeleter(sid, ctx)
			synctest.Wait()
			requireClosed(t, d.entered, "deleter did not enter")

			s := startSearcher(sid, ctx)
			i := startInserter(sid, ctx)
			synctest.Wait()
			requireOpen(t, s.entered, "searcher entered while deleter was active")
			requireOpen(t, i.entered, "inserter entered while deleter was active")

			close(d.release)
			synctest.Wait()
			requireSuccess(t, d.done, "deleter did not complete")
			requireClosed(t, s.entered, "searcher did not enter after deleter left")
			requireClosed(t, i.entered, "inserter did not enter after deleter left")

			close(s.release)
			close(i.release)
			synctest.Wait()
			requireSuccess(t, s.done, "searcher did not complete")
			requireSuccess(t, i.done, "inserter did not complete")
		})
	})

	t.Run("DeleterExclusive", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			ctx := t.Context()

			d1 := startDeleter(sid, ctx)
			synctest.Wait()
			requireClosed(t, d1.entered, "first deleter did not enter")

			d2 := startDeleter(sid, ctx)
			synctest.Wait()
			requireOpen(t, d2.entered, "second deleter entered while first deleter was active")

			close(d1.release)
			synctest.Wait()
			requireSuccess(t, d1.done, "first deleter did not complete")
			requireClosed(t, d2.entered, "second deleter did not enter after first deleter left")

			close(d2.release)
			synctest.Wait()
			requireSuccess(t, d2.done, "second deleter did not complete")
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			ctx := t.Context()
			runSearchersShare(t, sid, ctx)
			runInserterExclusive(t, sid, ctx)
			runSearchersShare(t, sid, ctx)
		})
	})

	t.Run("CancelWhileWaitingSearcher", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			active := startDeleter(sid, t.Context())
			synctest.Wait()
			requireClosed(t, active.entered, "deleter did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			searcher := startSearcher(sid, ctx)
			synctest.Wait()
			requireOpen(t, searcher.entered, "searcher entered while deleter was active")

			cancel()
			synctest.Wait()
			requireCanceled(t, searcher.done, "searcher did not exit after cancellation")

			close(active.release)
			synctest.Wait()
			requireSuccess(t, active.done, "deleter did not complete")
		})
	})

	t.Run("CancelWhileWaitingInserter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			active := startDeleter(sid, t.Context())
			synctest.Wait()
			requireClosed(t, active.entered, "deleter did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			inserter := startInserter(sid, ctx)
			synctest.Wait()
			requireOpen(t, inserter.entered, "inserter entered while deleter was active")

			cancel()
			synctest.Wait()
			requireCanceled(t, inserter.done, "inserter did not exit after cancellation")

			close(active.release)
			synctest.Wait()
			requireSuccess(t, active.done, "deleter did not complete")
		})
	})

	t.Run("CancelWhileWaitingDeleter", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sid := newImpl()
			activeSearcher := startSearcher(sid, t.Context())
			activeInserter := startInserter(sid, t.Context())
			synctest.Wait()
			requireClosed(t, activeSearcher.entered, "searcher did not enter")
			requireClosed(t, activeInserter.entered, "inserter did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			deleter := startDeleter(sid, ctx)
			synctest.Wait()
			requireOpen(t, deleter.entered, "deleter entered while searcher and inserter were active")

			cancel()
			synctest.Wait()
			requireCanceled(t, deleter.done, "deleter did not exit after cancellation")

			close(activeSearcher.release)
			close(activeInserter.release)
			synctest.Wait()
			requireSuccess(t, activeSearcher.done, "searcher did not complete")
			requireSuccess(t, activeInserter.done, "inserter did not complete")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() SearchInsertDelete) {
	b.Helper()
	b.ReportAllocs()

	b.Run("SearcherInserter", func(b *testing.B) {
		for b.Loop() {
			ctx := b.Context()
			sid := newImpl()
			done := make(chan error, 1)
			go func() {
				done <- sid.Searcher(ctx, func() {})
			}()
			if err := sid.Inserter(ctx, func() {}); err != nil {
				b.Fatalf("inserter returned %v", err)
			}
			if err := <-done; err != nil {
				b.Fatalf("searcher returned %v", err)
			}
		}
	})

	b.Run("Deleter", func(b *testing.B) {
		ctx := b.Context()
		sid := newImpl()
		for b.Loop() {
			if err := sid.Deleter(ctx, func() {}); err != nil {
				b.Fatalf("deleter returned %v", err)
			}
		}
	})
}
