//go:build challenge

// Package extendeddininghall contains the challenge version of the extended
// dining hall problem.
//
// Students first get food, then wait until it is socially safe to start dining,
// and finally wait again if leaving would strand exactly one diner alone.
package extendeddininghall

import (
	"context"
	"sync"
)

type ExtendedDiningHall struct {
	mu             sync.Mutex
	readyToEat     int
	eating         int
	readyToLeave   int
	waitingSit     chan struct{}
	waitingToLeave chan struct{}
}

func NewExtendedDiningHall() *ExtendedDiningHall {
	return &ExtendedDiningHall{}
}

func (h *ExtendedDiningHall) Student(ctx context.Context, getFood, dine, leave func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	getFood()

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := h.waitToSit(ctx); err != nil {
		return err
	}

	dine()

	if err := h.waitToLeave(ctx); err != nil {
		return err
	}

	leave()
	return nil
}

func (h *ExtendedDiningHall) waitToSit(ctx context.Context) error {
	h.mu.Lock()
	h.readyToEat++

	switch {
	case h.eating == 0 && h.readyToEat == 1:
		waiter := make(chan struct{})
		h.waitingSit = waiter
		h.mu.Unlock()

		select {
		case <-waiter:
			return nil
		case <-ctx.Done():
			h.mu.Lock()
			if h.waitingSit == waiter {
				h.waitingSit = nil
				h.readyToEat--
				h.mu.Unlock()
				return ctx.Err()
			}
			h.mu.Unlock()
			<-waiter
			return nil
		}
	case h.eating == 0 && h.readyToEat == 2:
		waiter := h.waitingSit
		h.waitingSit = nil
		h.readyToEat -= 2
		h.eating += 2
		h.mu.Unlock()
		if waiter != nil {
			close(waiter)
		}
		return nil
	default:
		h.readyToEat--
		h.eating++
		if h.eating == 2 && h.readyToLeave == 1 && h.waitingToLeave != nil {
			waiter := h.waitingToLeave
			h.waitingToLeave = nil
			h.readyToLeave--
			h.mu.Unlock()
			close(waiter)
			return nil
		}
		h.mu.Unlock()
		return nil
	}
}

func (h *ExtendedDiningHall) waitToLeave(ctx context.Context) error {
	h.mu.Lock()
	h.eating--
	h.readyToLeave++

	switch {
	case h.eating == 1 && h.readyToLeave == 1:
		waiter := make(chan struct{})
		h.waitingToLeave = waiter
		h.mu.Unlock()

		select {
		case <-waiter:
			return nil
		case <-ctx.Done():
			h.mu.Lock()
			if h.waitingToLeave == waiter {
				h.waitingToLeave = nil
				h.readyToLeave--
				h.mu.Unlock()
				return ctx.Err()
			}
			h.mu.Unlock()
			<-waiter
			return nil
		}
	case h.eating == 0 && h.readyToLeave == 2:
		waiter := h.waitingToLeave
		h.waitingToLeave = nil
		h.readyToLeave -= 2
		h.mu.Unlock()
		if waiter != nil {
			close(waiter)
		}
		return nil
	default:
		h.readyToLeave--
		h.mu.Unlock()
		return nil
	}
}
