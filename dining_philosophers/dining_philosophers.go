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

type DiningPhilosophers struct {
	forks [5]chan struct{}
}

func NewDiningPhilosophers() *DiningPhilosophers {
	forks := [5]chan struct{}{}
	for i := range forks {
		forks[i] = make(chan struct{}, 1)
		forks[i] <- struct{}{}
	}
	return &DiningPhilosophers{forks}
}

func (dp *DiningPhilosophers) Dine(ctx context.Context, i int, think, eat func()) {
	think()
	dp.getForks(i)
	eat()
	dp.putForks(i)
}

func (dp *DiningPhilosophers) getForks(i int) {
	<-dp.forks[left(i)]
	<-dp.forks[right(i)]
}

func (dp *DiningPhilosophers) putForks(i int) {
	dp.forks[left(i)] <- struct{}{}
	dp.forks[right(i)] <- struct{}{}
}

func left(i int) int {
	return i
}

func right(i int) int {
	return (i + 1) % 5
}
