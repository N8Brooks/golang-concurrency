package dining_philosophers

type DiningPhilosophers struct {
	forks [5]chan struct{}
}

func NewDiningPhilosophers() *DiningPhilosophers {
	dp := DiningPhilosophers{}
	dp.forks = [5]chan struct{}{
		make(chan struct{}, 1),
		make(chan struct{}, 1),
		make(chan struct{}, 1),
		make(chan struct{}, 1),
		make(chan struct{}, 1),
	}
	return &dp
}

func (dp *DiningPhilosophers) WantsToEat(philosopher int, pickLeftFork, pickRightFork, eat, putLeftFork, putRightFork func()) {
	if philosopher == 4 {
		dp.forks[0] <- struct{}{}
		pickRightFork()
		dp.forks[4] <- struct{}{}
		pickLeftFork()
	} else {
		dp.forks[philosopher] <- struct{}{}
		pickLeftFork()
		dp.forks[philosopher+1] <- struct{}{}
		pickRightFork()
	}

	eat()

	if philosopher == 4 {
		putRightFork()
		<-dp.forks[0]
		putLeftFork()
		<-dp.forks[4]
	} else {
		putLeftFork()
		<-dp.forks[philosopher]
		putRightFork()
		<-dp.forks[philosopher+1]
	}
}
