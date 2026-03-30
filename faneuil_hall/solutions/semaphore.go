// Package solutions contains implementations of the Faneuil Hall problem.
package solutions

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

var closedSignal = func() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()

type Semaphore struct {
	noJudge *semaphore.Weighted

	mu           sync.Mutex
	entered      int
	checked      int
	certified    int
	judgePresent bool
	allChecked   chan struct{}
	confirmed    chan struct{}
	allGone      chan struct{}
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		noJudge:   semaphore.NewWeighted(1),
		confirmed: make(chan struct{}),
		allGone:   closedSignal,
	}
}

func (fh *Semaphore) Immigrant(ctx context.Context, enter, checkIn, sitDown, swear, getCertificate, leave func()) error {
	if err := fh.noJudge.Acquire(ctx, 1); err != nil {
		return err
	}
	enter()

	fh.mu.Lock()
	fh.entered++
	confirmed := fh.confirmed
	fh.mu.Unlock()

	fh.noJudge.Release(1)

	checkIn()

	fh.mu.Lock()
	fh.checked++
	if fh.judgePresent && fh.entered == fh.checked && fh.allChecked != nil {
		close(fh.allChecked)
		fh.allChecked = nil
	}
	fh.mu.Unlock()

	sitDown()
	<-confirmed

	swear()
	getCertificate()

	if err := fh.noJudge.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	leave()
	fh.noJudge.Release(1)

	fh.mu.Lock()
	if fh.certified > 0 {
		fh.certified--
		if fh.certified == 0 {
			close(fh.allGone)
		}
	}
	fh.mu.Unlock()

	return nil
}

func (fh *Semaphore) Judge(ctx context.Context, enter, confirm, leave func()) error {
	if err := fh.waitForAllGone(ctx); err != nil {
		return err
	}
	if err := fh.noJudge.Acquire(ctx, 1); err != nil {
		return err
	}

	enter()

	fh.mu.Lock()
	fh.judgePresent = true
	if fh.entered > fh.checked {
		fh.allChecked = make(chan struct{})
		allChecked := fh.allChecked
		fh.mu.Unlock()
		<-allChecked
	} else {
		fh.mu.Unlock()
	}

	confirm()

	fh.mu.Lock()
	n := fh.checked
	fh.certified = n
	if n > 0 {
		fh.allGone = make(chan struct{})
	} else {
		fh.allGone = closedSignal
	}
	close(fh.confirmed)
	fh.confirmed = make(chan struct{})
	fh.entered = 0
	fh.checked = 0
	fh.mu.Unlock()

	leave()

	fh.mu.Lock()
	fh.judgePresent = false
	fh.mu.Unlock()

	fh.noJudge.Release(1)
	return nil
}

func (fh *Semaphore) Spectator(ctx context.Context, enter, spectate, leave func()) error {
	if err := fh.noJudge.Acquire(ctx, 1); err != nil {
		return err
	}
	enter()
	fh.noJudge.Release(1)

	spectate()
	leave()
	return nil
}

func (fh *Semaphore) waitForAllGone(ctx context.Context) error {
	fh.mu.Lock()
	if fh.certified == 0 {
		fh.mu.Unlock()
		return nil
	}
	allGone := fh.allGone
	fh.mu.Unlock()

	select {
	case <-allGone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
