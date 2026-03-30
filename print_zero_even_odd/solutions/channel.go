// Package solutions contains implementations of the print zero even odd problem.
package solutions

type Channel struct {
	n    int
	zero chan struct{}
	even chan struct{}
	odd  chan struct{}
}

func NewChannel(n int) *Channel {
	return &Channel{
		n:    n,
		zero: make(chan struct{}, 1),
		even: make(chan struct{}),
		odd:  make(chan struct{}),
	}
}

func (z *Channel) Zero(printNumber func(int)) {
	z.zero <- struct{}{}
	for i := 1; i <= z.n; i++ {
		<-z.zero
		printNumber(0)
		if i&1 == 1 {
			z.odd <- struct{}{}
		} else {
			z.even <- struct{}{}
		}
	}
}

func (z *Channel) Even(printNumber func(int)) {
	for i := 2; i <= z.n; i += 2 {
		<-z.even
		printNumber(i)
		z.zero <- struct{}{}
	}
}

func (z *Channel) Odd(printNumber func(int)) {
	for i := 1; i <= z.n; i += 2 {
		<-z.odd
		printNumber(i)
		z.zero <- struct{}{}
	}
}
