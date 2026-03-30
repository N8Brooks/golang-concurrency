//go:build challenge

// Package dininghall contains the challenge version of the dining hall
// problem.
//
// Students dine concurrently, but a student may not leave if doing so would
// leave exactly one other student still dining alone.
package dininghall

import (
	"context"
	"sync"
)

type DiningHall struct {
	mu           sync.Mutex
	eating       int
	readyToLeave int
	waitingLeave chan struct{}
}

func NewDiningHall() *DiningHall {
	return &DiningHall{}
}

func (h *DiningHall) Student(ctx context.Context, dine, leave func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	h.mu.Lock()
	h.eating++
	if h.eating == 2 && h.readyToLeave == 1 && h.waitingLeave != nil {
		waiter := h.waitingLeave
		h.waitingLeave = nil
		h.readyToLeave--
		close(waiter)
	}
	h.mu.Unlock()

	dine()

	h.mu.Lock()
	h.eating--
	h.readyToLeave++

	switch {
	case h.eating == 1 && h.readyToLeave == 1:
		waiter := make(chan struct{})
		h.waitingLeave = waiter
		h.mu.Unlock()

		select {
		case <-waiter:
		case <-ctx.Done():
			h.mu.Lock()
			if h.waitingLeave == waiter {
				h.waitingLeave = nil
				h.readyToLeave--
				h.mu.Unlock()
				return ctx.Err()
			}
			h.mu.Unlock()
			<-waiter
		}
	case h.eating == 0 && h.readyToLeave == 2:
		waiter := h.waitingLeave
		h.waitingLeave = nil
		h.readyToLeave -= 2
		h.mu.Unlock()
		if waiter != nil {
			close(waiter)
		}
	default:
		h.readyToLeave--
		h.mu.Unlock()
	}

	leave()
	return nil
}
