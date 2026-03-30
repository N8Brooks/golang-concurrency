// Package solutions contains implementations of the search-insert-delete problem.
package solutions

import (
	"context"
	"sync"
)

type Semaphore struct {
	insertMutex  *semaphore
	noSearcher   *semaphore
	noInserter   *semaphore
	searchSwitch *lightswitch
	insertSwitch *lightswitch
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		insertMutex:  newSemaphore(1),
		noSearcher:   newSemaphore(1),
		noInserter:   newSemaphore(1),
		searchSwitch: newLightswitch(),
		insertSwitch: newLightswitch(),
	}
}

func (s *Semaphore) Searcher(ctx context.Context, search func()) error {
	if err := s.searchSwitch.Lock(ctx, s.noSearcher); err != nil {
		return err
	}
	defer s.searchSwitch.Unlock(s.noSearcher)
	search()
	return nil
}

func (s *Semaphore) Inserter(ctx context.Context, insert func()) error {
	if err := s.insertSwitch.Lock(ctx, s.noInserter); err != nil {
		return err
	}
	defer s.insertSwitch.Unlock(s.noInserter)

	if err := s.insertMutex.Wait(ctx); err != nil {
		return err
	}
	defer s.insertMutex.Signal()
	insert()
	return nil
}

func (s *Semaphore) Deleter(ctx context.Context, remove func()) error {
	if err := s.noSearcher.Wait(ctx); err != nil {
		return err
	}

	if err := s.noInserter.Wait(ctx); err != nil {
		s.noSearcher.Signal()
		return err
	}
	defer s.noSearcher.Signal()
	defer s.noInserter.Signal()

	remove()
	return nil
}

type lightswitch struct {
	mu      sync.Mutex
	count   int
	waiting int
}

func newLightswitch() *lightswitch {
	return &lightswitch{}
}

func (l *lightswitch) Lock(ctx context.Context, sem *semaphore) error {
	l.mu.Lock()
	if l.count > 0 {
		l.count++
		l.mu.Unlock()
		return nil
	}
	l.waiting++
	l.mu.Unlock()

	if err := sem.Wait(ctx); err != nil {
		l.mu.Lock()
		l.waiting--
		l.mu.Unlock()
		return err
	}

	l.mu.Lock()
	l.waiting--
	l.count++
	l.mu.Unlock()
	return nil
}

func (l *lightswitch) Unlock(sem *semaphore) {
	l.mu.Lock()
	l.count--
	last := l.count == 0
	l.mu.Unlock()
	if last {
		sem.Signal()
	}
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
