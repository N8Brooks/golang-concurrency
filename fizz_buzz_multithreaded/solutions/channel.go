// Package solutions contains implementations of the multithreaded FizzBuzz
// problem.
package solutions

type Channel struct {
	n        int
	fizz     chan struct{}
	buzz     chan struct{}
	fizzbuzz chan struct{}
	ack      chan struct{}
}

func NewChannel(n int) *Channel {
	return &Channel{
		n:        n,
		fizz:     make(chan struct{}),
		buzz:     make(chan struct{}),
		fizzbuzz: make(chan struct{}),
		ack:      make(chan struct{}),
	}
}

func (fb *Channel) Fizz(printFizz func()) {
	for range fb.fizz {
		printFizz()
		fb.ack <- struct{}{}
	}
}

func (fb *Channel) Buzz(printBuzz func()) {
	for range fb.buzz {
		printBuzz()
		fb.ack <- struct{}{}
	}
}

func (fb *Channel) FizzBuzz(printFizzBuzz func()) {
	for range fb.fizzbuzz {
		printFizzBuzz()
		fb.ack <- struct{}{}
	}
}

func (fb *Channel) Number(printNumber func(int)) {
	defer close(fb.fizz)
	defer close(fb.buzz)
	defer close(fb.fizzbuzz)

	for i := 1; i <= fb.n; i++ {
		switch {
		case i%15 == 0:
			fb.fizzbuzz <- struct{}{}
			<-fb.ack
		case i%5 == 0:
			fb.buzz <- struct{}{}
			<-fb.ack
		case i%3 == 0:
			fb.fizz <- struct{}{}
			<-fb.ack
		default:
			printNumber(i)
		}
	}
}
