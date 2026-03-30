package solutions

import (
	"context"
	"sync"
)

type customerWaiter struct {
	ready        chan struct{}
	customerDone chan struct{}
	barberDone   chan struct{}
}

type barberWaiter struct {
	customer chan *customerWaiter
}

type Semaphore struct {
	mu               sync.Mutex
	capacity         int
	customers        int
	waitingCustomers []*customerWaiter
	sleepingBarbers  []*barberWaiter
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{capacity: capacity}
}

func (b *Semaphore) Customer(ctx context.Context, getHairCut, balk func()) error {
	waiter := &customerWaiter{
		ready:        make(chan struct{}),
		customerDone: make(chan struct{}),
		barberDone:   make(chan struct{}),
	}

	b.mu.Lock()
	if b.customers == b.capacity {
		b.mu.Unlock()
		balk()
		return nil
	}
	b.customers++
	if len(b.sleepingBarbers) > 0 {
		barber := b.sleepingBarbers[0]
		copy(b.sleepingBarbers, b.sleepingBarbers[1:])
		b.sleepingBarbers[len(b.sleepingBarbers)-1] = nil
		b.sleepingBarbers = b.sleepingBarbers[:len(b.sleepingBarbers)-1]
		b.mu.Unlock()
		barber.customer <- waiter
	} else {
		b.waitingCustomers = append(b.waitingCustomers, waiter)
		b.mu.Unlock()
	}

	select {
	case <-waiter.ready:
	case <-ctx.Done():
		b.mu.Lock()
		var removed bool
		b.waitingCustomers, removed = removeCustomerWaiter(b.waitingCustomers, waiter)
		if removed {
			b.customers--
			b.mu.Unlock()
			return ctx.Err()
		}
		b.mu.Unlock()
		<-waiter.ready
	}

	getHairCut()
	close(waiter.customerDone)
	<-waiter.barberDone

	b.mu.Lock()
	b.customers--
	b.mu.Unlock()
	return nil
}

func (b *Semaphore) Barber(ctx context.Context, cutHair func()) error {
	var waiter *customerWaiter

	b.mu.Lock()
	if len(b.waitingCustomers) > 0 {
		waiter = b.waitingCustomers[0]
		copy(b.waitingCustomers, b.waitingCustomers[1:])
		b.waitingCustomers[len(b.waitingCustomers)-1] = nil
		b.waitingCustomers = b.waitingCustomers[:len(b.waitingCustomers)-1]
		b.mu.Unlock()
	} else {
		barber := &barberWaiter{customer: make(chan *customerWaiter, 1)}
		b.sleepingBarbers = append(b.sleepingBarbers, barber)
		b.mu.Unlock()

		select {
		case waiter = <-barber.customer:
		case <-ctx.Done():
			b.mu.Lock()
			var removed bool
			b.sleepingBarbers, removed = removeBarberWaiter(b.sleepingBarbers, barber)
			b.mu.Unlock()
			if removed {
				return ctx.Err()
			}
			waiter = <-barber.customer
		}
	}

	close(waiter.ready)
	cutHair()
	<-waiter.customerDone
	close(waiter.barberDone)
	return nil
}

func removeCustomerWaiter(waiters []*customerWaiter, target *customerWaiter) ([]*customerWaiter, bool) {
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

func removeBarberWaiter(waiters []*barberWaiter, target *barberWaiter) ([]*barberWaiter, bool) {
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
