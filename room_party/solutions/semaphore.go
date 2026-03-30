package solutions

import (
	"context"
	"sync"
)

const breakupThreshold = 50

type deanState uint8

const (
	deanNotHere deanState = iota
	deanWaiting
	deanInRoom
)

type Semaphore struct {
	mu             sync.Mutex
	students       int
	dean           deanState
	deanWaiter     chan struct{}
	studentWaiters []chan struct{}
}

func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

func (r *Semaphore) Student(ctx context.Context, party func()) error {
	r.mu.Lock()
	for r.dean == deanInRoom {
		if err := r.waitStudentLocked(ctx); err != nil {
			r.mu.Unlock()
			return err
		}
	}
	r.students++
	if r.students >= breakupThreshold && r.dean == deanWaiting {
		r.signalDeanLocked()
	}
	r.mu.Unlock()

	party()

	r.mu.Lock()
	r.students--
	if r.students == 0 && (r.dean == deanWaiting || r.dean == deanInRoom) {
		r.signalDeanLocked()
	}
	r.mu.Unlock()
	return nil
}

func (r *Semaphore) Dean(ctx context.Context, search, breakup func()) error {
	r.mu.Lock()
	if r.students > 0 && r.students < breakupThreshold {
		r.dean = deanWaiting
		for r.students > 0 && r.students < breakupThreshold {
			if err := r.waitDeanLocked(ctx); err != nil {
				r.dean = deanNotHere
				r.mu.Unlock()
				return err
			}
		}
	}

	r.dean = deanInRoom

	if r.students >= breakupThreshold {
		r.mu.Unlock()
		breakup()
		r.mu.Lock()
		for r.students > 0 {
			if err := r.waitDeanLocked(context.Background()); err != nil {
				panic(err)
			}
		}
		r.dean = deanNotHere
		r.wakeStudentsLocked()
		r.mu.Unlock()
		return nil
	}

	r.mu.Unlock()
	search()
	r.mu.Lock()
	r.dean = deanNotHere
	r.wakeStudentsLocked()
	r.mu.Unlock()
	return nil
}

func (r *Semaphore) waitStudentLocked(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	waiter := make(chan struct{})
	r.studentWaiters = append(r.studentWaiters, waiter)
	r.mu.Unlock()

	select {
	case <-waiter:
		r.mu.Lock()
		return nil
	case <-ctx.Done():
		r.mu.Lock()
		if removeWaiter(&r.studentWaiters, waiter) {
			return ctx.Err()
		}
		return nil
	}
}

func (r *Semaphore) waitDeanLocked(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if r.deanWaiter != nil {
		panic("multiple Dean waiters")
	}

	waiter := make(chan struct{})
	r.deanWaiter = waiter
	r.mu.Unlock()

	select {
	case <-waiter:
		r.mu.Lock()
		return nil
	case <-ctx.Done():
		r.mu.Lock()
		if r.deanWaiter == waiter {
			r.deanWaiter = nil
			return ctx.Err()
		}
		return nil
	}
}

func (r *Semaphore) signalDeanLocked() {
	if r.deanWaiter == nil {
		return
	}
	close(r.deanWaiter)
	r.deanWaiter = nil
}

func (r *Semaphore) wakeStudentsLocked() {
	for _, waiter := range r.studentWaiters {
		close(waiter)
	}
	r.studentWaiters = nil
}

func removeWaiter(waiters *[]chan struct{}, target chan struct{}) bool {
	for i, waiter := range *waiters {
		if waiter != target {
			continue
		}
		copy((*waiters)[i:], (*waiters)[i+1:])
		(*waiters)[len(*waiters)-1] = nil
		*waiters = (*waiters)[:len(*waiters)-1]
		return true
	}
	return false
}
