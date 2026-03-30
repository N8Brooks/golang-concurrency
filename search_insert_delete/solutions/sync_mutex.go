// Package solutions contains implementations of the search-insert-delete problem.
package solutions

import (
	"context"
	"sync"
)

type SyncMutex struct {
	mu        sync.Mutex
	cond      *sync.Cond
	searchers int
	inserters int
	deleter   bool
}

func NewSyncMutex() *SyncMutex {
	s := &SyncMutex{}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *SyncMutex) Searcher(ctx context.Context, search func()) error {
	if err := s.wait(ctx, func() bool {
		return !s.deleter
	}, func() {
		s.searchers++
	}); err != nil {
		return err
	}

	defer func() {
		s.mu.Lock()
		s.searchers--
		s.mu.Unlock()
		s.cond.Broadcast()
	}()

	search()
	return nil
}

func (s *SyncMutex) Inserter(ctx context.Context, insert func()) error {
	if err := s.wait(ctx, func() bool {
		return !s.deleter && s.inserters == 0
	}, func() {
		s.inserters = 1
	}); err != nil {
		return err
	}

	defer func() {
		s.mu.Lock()
		s.inserters = 0
		s.mu.Unlock()
		s.cond.Broadcast()
	}()

	insert()
	return nil
}

func (s *SyncMutex) Deleter(ctx context.Context, remove func()) error {
	if err := s.wait(ctx, func() bool {
		return !s.deleter && s.searchers == 0 && s.inserters == 0
	}, func() {
		s.deleter = true
	}); err != nil {
		return err
	}

	defer func() {
		s.mu.Lock()
		s.deleter = false
		s.mu.Unlock()
		s.cond.Broadcast()
	}()

	remove()
	return nil
}

func (s *SyncMutex) wait(ctx context.Context, ready func() bool, enter func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			s.cond.Broadcast()
			s.mu.Unlock()
		case <-done:
		}
	}()

	for !ready() && ctx.Err() == nil {
		s.cond.Wait()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	enter()
	return nil
}
