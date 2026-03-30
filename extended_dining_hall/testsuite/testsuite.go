// Package testsuite contains reusable behavioral tests for extended dining hall implementations.
package testsuite

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
)

type ExtendedDiningHall interface {
	Student(ctx context.Context, getFood, dine, leave func()) error
}

type studentRun struct {
	gotFood     chan struct{}
	dining      chan struct{}
	releaseDine chan struct{}
	left        chan struct{}
	done        chan error
}

func startStudent(hall ExtendedDiningHall, ctx context.Context) studentRun {
	run := studentRun{
		gotFood:     make(chan struct{}),
		dining:      make(chan struct{}),
		releaseDine: make(chan struct{}),
		left:        make(chan struct{}),
		done:        make(chan error, 1),
	}

	go func() {
		run.done <- hall.Student(ctx, func() {
			close(run.gotFood)
		}, func() {
			close(run.dining)
			<-run.releaseDine
		}, func() {
			close(run.left)
		})
	}()

	return run
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

func requireSuccess(t *testing.T, done <-chan error, message string) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("%s: %v", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func requireCanceled(t *testing.T, done <-chan error, message string) {
	t.Helper()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: got %v, want context.Canceled", message, err)
		}
	default:
		t.Fatal(message)
	}
}

func runPair(t *testing.T, hall ExtendedDiningHall, ctx context.Context) {
	t.Helper()

	first := startStudent(hall, ctx)
	synctest.Wait()
	requireClosed(t, first.gotFood, "first student did not get food")
	requireOpen(t, first.dining, "first student started dining alone")

	second := startStudent(hall, ctx)
	synctest.Wait()
	requireClosed(t, second.gotFood, "second student did not get food")
	requireClosed(t, first.dining, "first student was not released by a second ready-to-eat student")
	requireClosed(t, second.dining, "second student did not start dining with the first student")

	close(first.releaseDine)
	close(second.releaseDine)
	synctest.Wait()

	requireClosed(t, first.left, "first student did not leave")
	requireClosed(t, second.left, "second student did not leave")
	requireSuccess(t, first.done, "first student did not complete")
	requireSuccess(t, second.done, "second student did not complete")
}

func Run(t *testing.T, newImpl func() ExtendedDiningHall) {
	t.Helper()

	t.Run("SecondArrivalReleasesWaitingStudent", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runPair(t, newImpl(), t.Context())
		})
	})

	t.Run("WaitingStudentReleasedByFinalDiner", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			hall := newImpl()
			ctx := t.Context()

			first := startStudent(hall, ctx)
			second := startStudent(hall, ctx)
			synctest.Wait()
			requireClosed(t, first.dining, "first student did not start dining")
			requireClosed(t, second.dining, "second student did not start dining")

			close(second.releaseDine)
			synctest.Wait()
			requireOpen(t, second.left, "second student left while first student was still dining alone")

			close(first.releaseDine)
			synctest.Wait()

			requireClosed(t, first.left, "first student did not leave after finishing")
			requireClosed(t, second.left, "waiting student was not released by the final diner")
			requireSuccess(t, first.done, "first student did not complete")
			requireSuccess(t, second.done, "second student did not complete")
		})
	})

	t.Run("WaitingStudentReleasedByArrival", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			hall := newImpl()
			ctx := t.Context()

			first := startStudent(hall, ctx)
			second := startStudent(hall, ctx)
			synctest.Wait()
			requireClosed(t, first.dining, "first student did not start dining")
			requireClosed(t, second.dining, "second student did not start dining")

			close(second.releaseDine)
			synctest.Wait()
			requireOpen(t, second.left, "second student left while first student was still dining alone")

			third := startStudent(hall, ctx)
			synctest.Wait()
			requireClosed(t, third.gotFood, "third student did not get food")
			requireClosed(t, third.dining, "third student did not start dining")
			requireClosed(t, second.left, "waiting student was not released by a new arrival")

			close(first.releaseDine)
			close(third.releaseDine)
			synctest.Wait()

			requireClosed(t, first.left, "first student did not leave")
			requireClosed(t, third.left, "third student did not leave")
			requireSuccess(t, first.done, "first student did not complete")
			requireSuccess(t, second.done, "second student did not complete")
			requireSuccess(t, third.done, "third student did not complete")
		})
	})

	t.Run("CancelWhileWaitingToSit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			hall := newImpl()
			ctx, cancel := context.WithCancel(t.Context())

			first := startStudent(hall, ctx)
			synctest.Wait()
			requireClosed(t, first.gotFood, "student did not get food")
			requireOpen(t, first.dining, "student started dining alone")

			cancel()
			synctest.Wait()
			requireCanceled(t, first.done, "student did not stop waiting to sit after cancellation")

			runPair(t, hall, t.Context())
		})
	})

	t.Run("CancelWhileWaitingToLeave", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			hall := newImpl()
			ctx := t.Context()

			first := startStudent(hall, ctx)
			secondCtx, cancel := context.WithCancel(ctx)
			second := startStudent(hall, secondCtx)
			synctest.Wait()
			requireClosed(t, first.dining, "first student did not start dining")
			requireClosed(t, second.dining, "second student did not start dining")

			close(second.releaseDine)
			synctest.Wait()
			requireOpen(t, second.left, "second student left while first student was still dining alone")

			cancel()
			synctest.Wait()
			requireCanceled(t, second.done, "waiting student did not stop waiting to leave after cancellation")

			close(first.releaseDine)
			synctest.Wait()
			requireClosed(t, first.left, "remaining student did not leave after finishing")
			requireSuccess(t, first.done, "remaining student did not complete")

			runPair(t, hall, ctx)
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			hall := newImpl()
			ctx := t.Context()
			runPair(t, hall, ctx)
			runPair(t, hall, ctx)
		})
	})
}

func Benchmark(b *testing.B, newImpl func() ExtendedDiningHall) {
	b.Helper()
	b.ReportAllocs()

	b.Run("PairedStudents", func(b *testing.B) {
		hall := newImpl()
		ctx := b.Context()
		for b.Loop() {
			first := startStudent(hall, ctx)
			second := startStudent(hall, ctx)
			<-first.dining
			<-second.dining
			close(first.releaseDine)
			close(second.releaseDine)
			if err := <-first.done; err != nil {
				b.Fatalf("first student returned %v", err)
			}
			if err := <-second.done; err != nil {
				b.Fatalf("second student returned %v", err)
			}
		}
	})
}
