// Package testsuite contains reusable behavioral tests for dining hall implementations.
package testsuite

import (
	"context"
	"testing"
	"testing/synctest"
)

type DiningHall interface {
	Student(ctx context.Context, dine, leave func()) error
}

type studentRun struct {
	dining      chan struct{}
	releaseDine chan struct{}
	left        chan struct{}
	done        chan error
}

func startStudent(hall DiningHall, ctx context.Context) studentRun {
	run := studentRun{
		dining:      make(chan struct{}),
		releaseDine: make(chan struct{}),
		left:        make(chan struct{}),
		done:        make(chan error, 1),
	}

	go func() {
		run.done <- hall.Student(ctx, func() {
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

func runOneShot(t *testing.T, hall DiningHall, ctx context.Context) {
	t.Helper()

	student := startStudent(hall, ctx)
	synctest.Wait()
	requireClosed(t, student.dining, "student did not start dining")

	close(student.releaseDine)
	synctest.Wait()
	requireClosed(t, student.left, "student did not leave after dining alone")
	requireSuccess(t, student.done, "student did not complete")
}

func Run(t *testing.T, newImpl func() DiningHall) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			runOneShot(t, newImpl(), t.Context())
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

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			hall := newImpl()
			ctx := t.Context()
			runOneShot(t, hall, ctx)

			first := startStudent(hall, ctx)
			second := startStudent(hall, ctx)
			synctest.Wait()
			close(second.releaseDine)
			synctest.Wait()
			requireOpen(t, second.left, "second student left too early")
			close(first.releaseDine)
			synctest.Wait()
			requireClosed(t, first.left, "first student did not leave")
			requireClosed(t, second.left, "second student did not leave")
			requireSuccess(t, first.done, "first student did not complete")
			requireSuccess(t, second.done, "second student did not complete")

			runOneShot(t, hall, ctx)
		})
	})
}

func Benchmark(b *testing.B, newImpl func() DiningHall) {
	b.Helper()
	b.ReportAllocs()

	b.Run("SingleStudent", func(b *testing.B) {
		hall := newImpl()
		ctx := b.Context()
		for b.Loop() {
			run := startStudent(hall, ctx)
			<-run.dining
			close(run.releaseDine)
			if err := <-run.done; err != nil {
				b.Fatalf("student returned %v", err)
			}
		}
	})

	b.Run("TwoStudents", func(b *testing.B) {
		ctx := b.Context()
		for b.Loop() {
			hall := newImpl()
			first := startStudent(hall, ctx)
			second := startStudent(hall, ctx)
			<-first.dining
			<-second.dining
			close(second.releaseDine)
			close(first.releaseDine)
			if err := <-first.done; err != nil {
				b.Fatalf("first student returned %v", err)
			}
			if err := <-second.done; err != nil {
				b.Fatalf("second student returned %v", err)
			}
		}
	})
}
