package solutions

type channelBarrier struct {
	n          int
	count      int
	mutex      chan struct{}
	turnstile1 chan struct{}
	turnstile2 chan struct{}
}

func newChannelBarrier(n int) *channelBarrier {
	mutex := make(chan struct{}, 1)
	mutex <- struct{}{}
	return &channelBarrier{
		n:          n,
		mutex:      mutex,
		turnstile1: make(chan struct{}, n),
		turnstile2: make(chan struct{}, n),
	}
}

func (b *channelBarrier) Wait() {
	<-b.mutex
	b.count++
	if b.count == b.n {
		for range b.n {
			b.turnstile1 <- struct{}{}
		}
	}
	b.mutex <- struct{}{}

	<-b.turnstile1

	<-b.mutex
	b.count--
	if b.count == 0 {
		for range b.n {
			b.turnstile2 <- struct{}{}
		}
	}
	b.mutex <- struct{}{}

	<-b.turnstile2
}

type Channel struct {
	mutex      chan struct{}
	oxygen     int
	hydrogen   int
	barrier    *channelBarrier
	oxyQueue   chan struct{}
	hydroQueue chan struct{}
}

func NewChannel() *Channel {
	mutex := make(chan struct{}, 1)
	mutex <- struct{}{}
	return &Channel{
		mutex:      mutex,
		barrier:    newChannelBarrier(3),
		oxyQueue:   make(chan struct{}, 1),
		hydroQueue: make(chan struct{}, 2),
	}
}

func (h2o *Channel) Hydrogen(bond func()) {
	<-h2o.mutex
	h2o.hydrogen++
	if h2o.hydrogen >= 2 && h2o.oxygen >= 1 {
		h2o.hydroQueue <- struct{}{}
		h2o.hydroQueue <- struct{}{}
		h2o.hydrogen -= 2
		h2o.oxyQueue <- struct{}{}
		h2o.oxygen--
	} else {
		h2o.mutex <- struct{}{}
	}

	<-h2o.hydroQueue
	bond()
	h2o.barrier.Wait()
}

func (h2o *Channel) Oxygen(bond func()) {
	<-h2o.mutex
	h2o.oxygen++
	if h2o.hydrogen >= 2 {
		h2o.hydroQueue <- struct{}{}
		h2o.hydroQueue <- struct{}{}
		h2o.hydrogen -= 2
		h2o.oxyQueue <- struct{}{}
		h2o.oxygen--
	} else {
		h2o.mutex <- struct{}{}
	}

	<-h2o.oxyQueue
	bond()
	h2o.barrier.Wait()
	h2o.mutex <- struct{}{}
}
