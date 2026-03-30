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
	mu            sync.Mutex
	cars          int
	capacity      int
	waiting       []*passengerWaiter
	loading       *carWaiter
	loadingArea   []*semaphore
	unloadingArea []*semaphore
}

func NewSemaphore(cars, capacity int) *Semaphore {
	r := &Semaphore{
		cars:          cars,
		capacity:      capacity,
		loadingArea:   make([]*semaphore, cars),
		unloadingArea: make([]*semaphore, cars),
	}
	for i := range cars {
		r.loadingArea[i] = newSemaphore(0)
		r.unloadingArea[i] = newSemaphore(0)
	}
	r.loadingArea[0].Signal()
	r.unloadingArea[0].Signal()
	return r
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

func (r *Semaphore) Car(ctx context.Context, car int, load, run, unload func()) error {
	if err := r.loadingArea[car].Wait(ctx); err != nil {
		return err
	}

	load()

	passengers, err := r.awaitPassengers(ctx)
	if err != nil {
		r.loadingArea[r.next(car)].Signal()
		return err
	}

	r.loadingArea[r.next(car)].Signal()

	for _, passenger := range passengers {
		close(passenger.boardAllowed)
	}
	for _, passenger := range passengers {
		<-passenger.boarded
	}

	run()

	if err := r.unloadingArea[car].Wait(context.Background()); err != nil {
		panic(err)
	}

	unload()

	for _, passenger := range passengers {
		close(passenger.unloadAllowed)
	}
	for _, passenger := range passengers {
		<-passenger.unboarded
	}

	r.unloadingArea[r.next(car)].Signal()
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

func (r *Semaphore) next(car int) int {
	return (car + 1) % r.cars
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

type semaphore struct {
	mu      sync.Mutex
	count   int
	waiters []chan struct{}
}

func newSemaphore(count int) *semaphore {
	return &semaphore{count: count}
}

func (s *semaphore) Signal() {
	s.mu.Lock()
	if len(s.waiters) > 0 {
		waiter := s.waiters[0]
		copy(s.waiters, s.waiters[1:])
		s.waiters[len(s.waiters)-1] = nil
		s.waiters = s.waiters[:len(s.waiters)-1]
		s.mu.Unlock()
		close(waiter)
		return
	}
	s.count++
	s.mu.Unlock()
}

func (s *semaphore) Wait(ctx context.Context) error {
	s.mu.Lock()
	if s.count > 0 {
		s.count--
		s.mu.Unlock()
		return nil
	}
	if err := ctx.Err(); err != nil {
		s.mu.Unlock()
		return err
	}

	waiter := make(chan struct{})
	s.waiters = append(s.waiters, waiter)
	s.mu.Unlock()

	select {
	case <-waiter:
		return nil
	case <-ctx.Done():
		s.mu.Lock()
		index := -1
		for i, candidate := range s.waiters {
			if candidate == waiter {
				index = i
				break
			}
		}
		if index >= 0 {
			copy(s.waiters[index:], s.waiters[index+1:])
			s.waiters[len(s.waiters)-1] = nil
			s.waiters = s.waiters[:len(s.waiters)-1]
			s.mu.Unlock()
			return ctx.Err()
		}
		s.mu.Unlock()
		<-waiter
		return nil
	}
}
