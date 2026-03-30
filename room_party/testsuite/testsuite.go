// Package testsuite contains reusable behavioral tests for room party implementations.
package testsuite

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
)

const BreakupThreshold = 50

type RoomParty interface {
	Student(ctx context.Context, party func()) error
	Dean(ctx context.Context, search, breakup func()) error
}

type studentRun struct {
	entered chan struct{}
	release chan struct{}
	done    chan error
}

type deanRun struct {
	searched       chan struct{}
	brokeUp        chan struct{}
	releaseSearch  chan struct{}
	releaseBreakup chan struct{}
	done           chan error
}

func startStudent(rp RoomParty, ctx context.Context) studentRun {
	run := studentRun{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		done:    make(chan error, 1),
	}

	go func() {
		run.done <- rp.Student(ctx, func() {
			close(run.entered)
			<-run.release
		})
	}()

	return run
}

func startDean(rp RoomParty, ctx context.Context) deanRun {
	run := deanRun{
		searched:       make(chan struct{}),
		brokeUp:        make(chan struct{}),
		releaseSearch:  make(chan struct{}),
		releaseBreakup: make(chan struct{}),
		done:           make(chan error, 1),
	}

	go func() {
		run.done <- rp.Dean(ctx, func() {
			close(run.searched)
			<-run.releaseSearch
		}, func() {
			close(run.brokeUp)
			<-run.releaseBreakup
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

func requirePending(t *testing.T, done <-chan error, message string) {
	t.Helper()

	select {
	case err := <-done:
		t.Fatalf("%s: returned %v", message, err)
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

func startStudents(rp RoomParty, ctx context.Context, count int) []studentRun {
	students := make([]studentRun, count)
	for i := range count {
		students[i] = startStudent(rp, ctx)
		synctest.Wait()
	}
	return students
}

func releaseStudents(students []studentRun) {
	for _, student := range students {
		close(student.release)
	}
}

func Run(t *testing.T, newImpl func() RoomParty) {
	t.Helper()

	t.Run("StudentsCanShareRoom", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()
			ctx := t.Context()

			first := startStudent(rp, ctx)
			synctest.Wait()
			requireClosed(t, first.entered, "first student did not enter")

			second := startStudent(rp, ctx)
			synctest.Wait()
			requireClosed(t, second.entered, "second student did not enter")

			close(first.release)
			close(second.release)
			synctest.Wait()

			requireSuccess(t, first.done, "first student did not complete")
			requireSuccess(t, second.done, "second student did not complete")
		})
	})

	t.Run("DeanSearchesEmptyRoom", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()
			ctx := t.Context()

			dean := startDean(rp, ctx)
			synctest.Wait()

			requireClosed(t, dean.searched, "Dean did not search an empty room")
			requireOpen(t, dean.brokeUp, "Dean broke up a party in an empty room")

			student := startStudent(rp, ctx)
			synctest.Wait()

			requireOpen(t, student.entered, "student entered while the Dean was searching")
			requirePending(t, student.done, "student completed while the Dean was searching")

			close(dean.releaseSearch)
			synctest.Wait()

			requireSuccess(t, dean.done, "Dean did not complete after searching")
			requireClosed(t, student.entered, "student did not enter after the Dean left")

			close(student.release)
			synctest.Wait()

			requireSuccess(t, student.done, "student did not complete")
		})
	})

	t.Run("DeanWaitsBelowThresholdThenSearches", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()
			ctx := t.Context()

			first := startStudent(rp, ctx)
			synctest.Wait()
			requireClosed(t, first.entered, "first student did not enter")

			dean := startDean(rp, ctx)
			synctest.Wait()

			requireOpen(t, dean.searched, "Dean searched while a student was still in the room")
			requireOpen(t, dean.brokeUp, "Dean broke up a party below the threshold")
			requirePending(t, dean.done, "Dean completed before the room became eligible")

			second := startStudent(rp, ctx)
			synctest.Wait()
			requireClosed(t, second.entered, "second student did not enter while the Dean was only waiting")

			close(first.release)
			close(second.release)
			synctest.Wait()

			requireClosed(t, dean.searched, "Dean did not search after the room became empty")
			requireOpen(t, dean.brokeUp, "Dean broke up the room instead of searching it")

			close(dean.releaseSearch)
			synctest.Wait()

			requireSuccess(t, first.done, "first student did not complete")
			requireSuccess(t, second.done, "second student did not complete")
			requireSuccess(t, dean.done, "Dean did not complete")
		})
	})

	t.Run("DeanBreaksUpPartyAndBlocksNewStudents", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()
			ctx := t.Context()

			students := startStudents(rp, ctx, BreakupThreshold-1)
			for i, student := range students {
				requireClosed(t, student.entered, fmt.Sprintf("student %d did not enter", i))
			}

			dean := startDean(rp, ctx)
			synctest.Wait()

			requireOpen(t, dean.searched, "Dean searched instead of waiting below the threshold")
			requireOpen(t, dean.brokeUp, "Dean broke up the party before the threshold was reached")

			thresholdStudent := startStudent(rp, ctx)
			synctest.Wait()

			requireClosed(t, thresholdStudent.entered, "threshold student did not enter")
			requireClosed(t, dean.brokeUp, "Dean did not break up the party at the threshold")
			requireOpen(t, dean.searched, "Dean searched instead of breaking up the party")

			blockedStudent := startStudent(rp, ctx)
			synctest.Wait()

			requireOpen(t, blockedStudent.entered, "student entered while the Dean was in the room")
			requirePending(t, blockedStudent.done, "blocked student completed while the Dean was in the room")

			close(students[0].release)
			synctest.Wait()

			requireSuccess(t, students[0].done, "student could not leave while the Dean was in the room")
			requirePending(t, dean.done, "Dean left before the room was empty")

			close(dean.releaseBreakup)
			synctest.Wait()

			requireOpen(t, blockedStudent.entered, "student entered before the Dean actually left")
			requirePending(t, dean.done, "Dean left before the remaining students were gone")

			releaseStudents(students[1:])
			close(thresholdStudent.release)
			synctest.Wait()

			requireSuccess(t, thresholdStudent.done, "threshold student did not complete")
			requireSuccess(t, dean.done, "Dean did not wait for the room to clear")
			requireClosed(t, blockedStudent.entered, "blocked student did not enter after the Dean left")

			close(blockedStudent.release)
			synctest.Wait()

			requireSuccess(t, blockedStudent.done, "blocked student did not complete")
		})
	})

	t.Run("CancelWaitingStudent", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()

			dean := startDean(rp, t.Context())
			synctest.Wait()
			requireClosed(t, dean.searched, "Dean did not begin searching")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			student := startStudent(rp, ctx)
			synctest.Wait()

			requireOpen(t, student.entered, "student entered while the Dean was searching")
			requirePending(t, student.done, "student completed before cancellation")

			cancel()
			synctest.Wait()

			select {
			case err := <-student.done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("student returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("student did not exit after cancellation")
			}

			requireOpen(t, student.entered, "student entered despite being canceled while waiting")

			close(dean.releaseSearch)
			synctest.Wait()

			requireSuccess(t, dean.done, "Dean did not complete")
		})
	})

	t.Run("CancelWaitingDean", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()

			student := startStudent(rp, t.Context())
			synctest.Wait()
			requireClosed(t, student.entered, "student did not enter")

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			dean := startDean(rp, ctx)
			synctest.Wait()

			requireOpen(t, dean.searched, "Dean searched while a student was inside")
			requireOpen(t, dean.brokeUp, "Dean broke up the room below the threshold")
			requirePending(t, dean.done, "Dean completed before cancellation")

			cancel()
			synctest.Wait()

			select {
			case err := <-dean.done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("Dean returned %v, want context.Canceled", err)
				}
			default:
				t.Fatal("Dean did not exit after cancellation")
			}

			extra := startStudent(rp, t.Context())
			synctest.Wait()
			requireClosed(t, extra.entered, "extra student did not enter after the waiting Dean canceled")

			close(student.release)
			close(extra.release)
			synctest.Wait()

			requireSuccess(t, student.done, "student did not complete")
			requireSuccess(t, extra.done, "extra student did not complete")
		})
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			rp := newImpl()
			ctx := t.Context()

			for range 3 {
				dean := startDean(rp, ctx)
				synctest.Wait()
				requireClosed(t, dean.searched, "Dean did not search an empty room")
				close(dean.releaseSearch)
				synctest.Wait()
				requireSuccess(t, dean.done, "Dean did not complete empty-room search")
			}

			students := startStudents(rp, ctx, BreakupThreshold)
			dean := startDean(rp, ctx)
			synctest.Wait()
			requireClosed(t, dean.brokeUp, "Dean did not break up a reusable party")
			close(dean.releaseBreakup)
			releaseStudents(students)
			synctest.Wait()

			for i, student := range students {
				requireSuccess(t, student.done, fmt.Sprintf("student %d did not complete reusable party", i))
			}
			requireSuccess(t, dean.done, "Dean did not complete reusable breakup")
		})
	})
}

func Benchmark(b *testing.B, newImpl func() RoomParty) {
	b.Helper()
	b.ReportAllocs()

	b.Run("StudentsOnly", func(b *testing.B) {
		rp := newImpl()
		ctx := b.Context()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := rp.Student(ctx, func() {}); err != nil {
					b.Fatalf("Student returned %v", err)
				}
			}
		})
	})

	b.Run("SearchEmpty", func(b *testing.B) {
		rp := newImpl()
		ctx := b.Context()

		for b.Loop() {
			if err := rp.Dean(ctx, func() {}, func() {}); err != nil {
				b.Fatalf("Dean returned %v", err)
			}
		}
	})
}
