package solutions

import "context"

const numPhilosophers = 5

type Footman struct {
	footman *semaphore
	forks   [numPhilosophers]*semaphore
}

func NewFootman() *Footman {
	dp := &Footman{
		footman: newSemaphore(numPhilosophers - 1),
	}
	for i := range numPhilosophers {
		dp.forks[i] = newSemaphore(1)
	}
	return dp
}

func (dp *Footman) Dine(ctx context.Context, philosopher int, think, eat func()) error {
	think()

	if err := dp.footman.Wait(ctx); err != nil {
		return err
	}

	right := philosopher
	left := (philosopher + 1) % numPhilosophers

	if err := dp.forks[right].Wait(ctx); err != nil {
		dp.footman.Signal()
		return err
	}
	if err := dp.forks[left].Wait(ctx); err != nil {
		dp.forks[right].Signal()
		dp.footman.Signal()
		return err
	}

	eat()

	dp.forks[left].Signal()
	dp.forks[right].Signal()
	dp.footman.Signal()
	return nil
}
