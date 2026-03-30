// Package solutions contains implementations of the extended child care
// problem.
package solutions

import (
	"context"
	"sync"
)

type SyncCond struct {
	mu       sync.Mutex
	cond     *sync.Cond
	adults   int
	children int
}

func NewSyncCond() *SyncCond {
	c := &SyncCond{}
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *SyncCond) Child(ctx context.Context, child func()) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	stop := context.AfterFunc(ctx, c.cond.Broadcast)
	defer stop()

	for c.children+1 > 3*c.adults {
		if err := ctx.Err(); err != nil {
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
	defer c.mu.Unlock()

	for c.children > 3*(c.adults-1) {
		c.cond.Wait()
	}
	c.adults--
	c.cond.Broadcast()
	return nil
}
