package solutions

import "context"

type philosopherState uint8

const (
	thinking philosopherState = iota
	hungry
	eating
)

type Tanenbaum struct {
	state [5]philosopherState
	ready [5]chan struct{}
	mu    chan struct{}
}

func NewTanenbaum() *Tanenbaum {
	dp := &Tanenbaum{
		mu: make(chan struct{}, 1),
	}
	dp.mu <- struct{}{}
	for i := range dp.ready {
		dp.ready[i] = make(chan struct{}, 1)
	}
	return dp
}

func (dp *Tanenbaum) Dine(ctx context.Context, philosopher int, think, eat func()) {
	think()

	if !dp.pickupForks(ctx, philosopher) {
		return
	}

	eat()
	dp.putDownForks(philosopher)
}

func (dp *Tanenbaum) pickupForks(ctx context.Context, philosopher int) bool {
	dp.lock()
	dp.state[philosopher] = hungry
	dp.test(philosopher)
	dp.unlock()

	select {
	case <-dp.ready[philosopher]:
		return true
	case <-ctx.Done():
		dp.lock()
		if dp.state[philosopher] != thinking {
			dp.state[philosopher] = thinking
			dp.test(leftPhilosopher(philosopher))
			dp.test(rightPhilosopher(philosopher))
		}
		dp.unlock()

		select {
		case <-dp.ready[philosopher]:
		default:
		}
		return false
	}
}

func (dp *Tanenbaum) putDownForks(philosopher int) {
	dp.lock()
	dp.state[philosopher] = thinking
	dp.test(leftPhilosopher(philosopher))
	dp.test(rightPhilosopher(philosopher))
	dp.unlock()
}

func (dp *Tanenbaum) test(philosopher int) {
	if dp.state[philosopher] != hungry {
		return
	}
	if dp.state[leftPhilosopher(philosopher)] == eating || dp.state[rightPhilosopher(philosopher)] == eating {
		return
	}

	dp.state[philosopher] = eating
	select {
	case dp.ready[philosopher] <- struct{}{}:
	default:
	}
}

func (dp *Tanenbaum) lock() {
	<-dp.mu
}

func (dp *Tanenbaum) unlock() {
	dp.mu <- struct{}{}
}

func leftPhilosopher(i int) int {
	return (i + 4) % 5
}

func rightPhilosopher(i int) int {
	return (i + 1) % 5
}
