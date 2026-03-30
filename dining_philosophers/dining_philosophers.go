//go:build challenge

// Package dining_philosophers contains the challenge version of the dining
// philosophers problem.
//
// Five philosophers sit around a table with five forks. Philosopher i needs
// fork i on the right and fork (i+1)%5 on the left before eating. Neighboring
// philosophers therefore contend for a shared fork, while non-neighbors should
// still be able to eat concurrently.
package dining_philosophers

import "context"

type DiningPhilosophers struct{}

func NewDiningPhilosophers() *DiningPhilosophers {
	return &DiningPhilosophers{}
}

func (dp *DiningPhilosophers) Dine(ctx context.Context, philosopher int, think, eat func()) error {
	panic("unimplemented")
}
