// Package solutions contains implementations of Hilzer's barbershop problem.
package solutions

import (
	"context"
	"sync"
)

type service struct {
	startCut    chan struct{}
	haircutDone chan struct{}
	paymentDone chan struct{}
	receiptDone chan struct{}
}

type customer struct {
	sofaReady  chan struct{}
	chairReady chan *service
}

type SyncMutex struct {
	mu            sync.Mutex
	customerReady *sync.Cond
	register      sync.Mutex

	inside    int
	sofaCount int
	standing  []*customer
	sofa      []*customer
}

func NewSyncMutex() *SyncMutex {
	s := &SyncMutex{}
	s.customerReady = sync.NewCond(&s.mu)
	return s
}

func (s *SyncMutex) Customer(ctx context.Context, enterShop, sitOnSofa, getHairCut, pay, exitShop, balk func()) error {
	c := &customer{
		sofaReady:  make(chan struct{}),
		chairReady: make(chan *service, 1),
	}

	s.mu.Lock()
	if s.inside == 20 {
		s.mu.Unlock()
		balk()
		return nil
	}
	s.inside++
	waitStanding := s.sofaCount == 4
	if waitStanding {
		s.standing = append(s.standing, c)
	} else {
		s.sofaCount++
	}
	s.mu.Unlock()

	enterShop()

	if waitStanding {
		select {
		case <-c.sofaReady:
		case <-ctx.Done():
			s.mu.Lock()
			var removed bool
			s.standing, removed = removeCustomer(s.standing, c)
			if removed {
				s.inside--
			}
			s.mu.Unlock()
			if removed {
				return ctx.Err()
			}
			<-c.sofaReady
		}
	}

	sitOnSofa()

	s.mu.Lock()
	s.sofa = append(s.sofa, c)
	s.customerReady.Signal()
	s.mu.Unlock()

	var srv *service
	select {
	case srv = <-c.chairReady:
	case <-ctx.Done():
		s.mu.Lock()
		var removed bool
		s.sofa, removed = removeCustomer(s.sofa, c)
		if removed {
			s.sofaCount--
			s.inside--
			s.promoteStandingLocked()
		}
		s.mu.Unlock()
		if removed {
			return ctx.Err()
		}
		srv = <-c.chairReady
	}

	<-srv.startCut
	getHairCut()
	<-srv.haircutDone
	pay()
	close(srv.paymentDone)
	<-srv.receiptDone
	exitShop()

	s.mu.Lock()
	s.inside--
	s.mu.Unlock()
	return nil
}

func (s *SyncMutex) Barber(ctx context.Context, cutHair, acceptPayment func()) error {
	stop := context.AfterFunc(ctx, func() {
		s.mu.Lock()
		s.customerReady.Broadcast()
		s.mu.Unlock()
	})
	defer stop()

	for {
		s.mu.Lock()
		for len(s.sofa) == 0 {
			if err := ctx.Err(); err != nil {
				s.mu.Unlock()
				return err
			}
			s.customerReady.Wait()
		}

		c := s.sofa[0]
		s.sofa = s.sofa[1:]
		s.sofaCount--
		s.promoteStandingLocked()

		srv := &service{
			startCut:    make(chan struct{}),
			haircutDone: make(chan struct{}),
			paymentDone: make(chan struct{}),
			receiptDone: make(chan struct{}),
		}
		s.mu.Unlock()

		c.chairReady <- srv
		close(srv.startCut)
		cutHair()
		close(srv.haircutDone)

		select {
		case <-srv.paymentDone:
		case <-ctx.Done():
			return ctx.Err()
		}

		s.register.Lock()
		acceptPayment()
		s.register.Unlock()
		close(srv.receiptDone)
	}
}

func (s *SyncMutex) promoteStandingLocked() {
	if len(s.standing) == 0 {
		return
	}
	next := s.standing[0]
	s.standing = s.standing[1:]
	s.sofaCount++
	close(next.sofaReady)
}

func removeCustomer(queue []*customer, target *customer) ([]*customer, bool) {
	for i, customer := range queue {
		if customer != target {
			continue
		}
		copy(queue[i:], queue[i+1:])
		queue[len(queue)-1] = nil
		return queue[:len(queue)-1], true
	}
	return queue, false
}
