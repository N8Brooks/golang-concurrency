package solutions

import (
	"context"
	"sync"
)

type passengerWaiter struct {
	boardAllowed  chan struct{}
	boarded       chan struct{}
	unloadAllowed chan struct{}
	unboarded     chan struct{}
}

type carWaiter struct {
	passengers chan []*passengerWaiter
}

type Semaphore struct {
	mu       sync.Mutex
	capacity int
	waiting  []*passengerWaiter
	loading  *carWaiter
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{capacity: capacity}
}

func (r *Semaphore) Passenger(ctx context.Context, board, unboard func()) error {
	waiter := &passengerWaiter{
		boardAllowed:  make(chan struct{}),
		boarded:       make(chan struct{}),
		unloadAllowed: make(chan struct{}),
		unboarded:     make(chan struct{}),
	}

	r.mu.Lock()
	r.waiting = append(r.waiting, waiter)
	car, passengers := r.tryAssembleRideLocked()
	r.mu.Unlock()

	if car != nil {
		car.passengers <- passengers
	}

	select {
	case <-waiter.boardAllowed:
	case <-ctx.Done():
		r.mu.Lock()
		var removed bool
		r.waiting, removed = removePassengerWaiter(r.waiting, waiter)
		r.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		<-waiter.boardAllowed
	}

	board()
	close(waiter.boarded)

	<-waiter.unloadAllowed
	unboard()
	close(waiter.unboarded)
	return nil
}

func (r *Semaphore) Car(ctx context.Context, load, run, unload func()) error {
	load()

	passengers, err := r.awaitPassengers(ctx)
	if err != nil {
		return err
	}

	for _, passenger := range passengers {
		close(passenger.boardAllowed)
	}
	for _, passenger := range passengers {
		<-passenger.boarded
	}

	run()
	unload()

	for _, passenger := range passengers {
		close(passenger.unloadAllowed)
	}
	for _, passenger := range passengers {
		<-passenger.unboarded
	}

	return nil
}

func (r *Semaphore) awaitPassengers(ctx context.Context) ([]*passengerWaiter, error) {
	r.mu.Lock()
	if len(r.waiting) >= r.capacity {
		passengers := append([]*passengerWaiter(nil), r.waiting[:r.capacity]...)
		r.waiting = r.waiting[r.capacity:]
		r.mu.Unlock()
		return passengers, nil
	}

	car := &carWaiter{passengers: make(chan []*passengerWaiter, 1)}
	r.loading = car
	r.mu.Unlock()

	select {
	case passengers := <-car.passengers:
		return passengers, nil
	case <-ctx.Done():
		r.mu.Lock()
		if r.loading == car {
			r.loading = nil
			r.mu.Unlock()
			return nil, ctx.Err()
		}
		r.mu.Unlock()
		return <-car.passengers, nil
	}
}

func (r *Semaphore) tryAssembleRideLocked() (*carWaiter, []*passengerWaiter) {
	if r.loading == nil || len(r.waiting) < r.capacity {
		return nil, nil
	}

	car := r.loading
	r.loading = nil

	passengers := append([]*passengerWaiter(nil), r.waiting[:r.capacity]...)
	r.waiting = r.waiting[r.capacity:]
	return car, passengers
}

func removePassengerWaiter(waiters []*passengerWaiter, target *passengerWaiter) ([]*passengerWaiter, bool) {
	for i, waiter := range waiters {
		if waiter != target {
			continue
		}
		copy(waiters[i:], waiters[i+1:])
		waiters[len(waiters)-1] = nil
		return waiters[:len(waiters)-1], true
	}
	return waiters, false
}
