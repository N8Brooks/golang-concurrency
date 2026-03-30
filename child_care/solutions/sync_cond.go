// Package solutions contains implementations of the child care problem.
package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	mu              sync.Mutex
	cond            *sync.Cond
	adults          int
	children        int
	departingAdults int
}

func NewSyncCond() *SyncCond {
	c := &SyncCond{}
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *SyncCond) Child(ctx context.Context, child func()) error {
	c.mu.Lock()
	stop := c.notifyOnCancel(ctx)
	defer stop()

	for c.children+1 > 3*(c.adults-c.departingAdults) {
		if err := ctx.Err(); err != nil {
			c.mu.Unlock()
			return err
		}
		c.cond.Wait()
	}
	c.children++
	c.mu.Unlock()

	child()

	c.mu.Lock()
	c.children--
	c.cond.Broadcast()
	c.mu.Unlock()
	return nil
}

func (c *SyncCond) Adult(ctx context.Context, adult func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	c.adults++
	c.cond.Broadcast()
	c.mu.Unlock()

	adult()

	c.mu.Lock()
	c.departingAdults++
	for c.children > 3*(c.adults-c.departingAdults) {
		c.cond.Wait()
	}
	c.departingAdults--
	c.adults--
	c.cond.Broadcast()
	c.mu.Unlock()
	return nil
}

func (c *SyncCond) notifyOnCancel(ctx context.Context) func() {
	if ctx.Done() == nil {
		return func() {}
	}

	stop := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.mu.Lock()
			c.cond.Broadcast()
			c.mu.Unlock()
		case <-stop:
		}
	}()
	return func() {
		close(stop)
	}
}
