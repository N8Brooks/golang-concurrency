package solutions

import "context"

type Asymmetric struct {
	forks [numPhilosophers]*semaphore
}

func NewAsymmetric() *Asymmetric {
	dp := &Asymmetric{}
	for i := range numPhilosophers {
		dp.forks[i] = newSemaphore(1)
	}
	return dp
}

func (dp *Asymmetric) Dine(ctx context.Context, philosopher int, think, eat func()) error {
	think()

	first := philosopher
	second := (philosopher + 1) % numPhilosophers
	if philosopher == 0 {
		first, second = second, first
	}

	if err := dp.forks[first].Wait(ctx); err != nil {
		return err
	}
	if err := dp.forks[second].Wait(ctx); err != nil {
		dp.forks[first].Signal()
		return err
	}

	eat()

	dp.forks[second].Signal()
	dp.forks[first].Signal()
	return nil
}
