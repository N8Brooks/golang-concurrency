package fizzbuzz

type FizzBuzz struct {
	n        int
	fizz     chan struct{}
	buzz     chan struct{}
	fizzbuzz chan struct{}
	number   chan struct{}
}

func NewFizzBuzz(n int) *FizzBuzz {
	fb := &FizzBuzz{n: n}
	fb.fizz = make(chan struct{})
	fb.buzz = make(chan struct{})
	fb.fizzbuzz = make(chan struct{})
	fb.number = make(chan struct{})
	return fb
}

func (fb *FizzBuzz) Fizz(printFizz func()) {
	for range fb.fizz {
		printFizz()
		fb.number <- struct{}{}
	}
}

func (fb *FizzBuzz) Buzz(printBuzz func()) {
	for range fb.buzz {
		printBuzz()
		fb.number <- struct{}{}
	}
}

func (fb *FizzBuzz) FizzBuzz(printFizzBuzz func()) {
	for range fb.fizzbuzz {
		printFizzBuzz()
		fb.number <- struct{}{}
	}
}

func (fb *FizzBuzz) Number(printNumber func(int)) {
	defer close(fb.fizz)
	defer close(fb.buzz)
	defer close(fb.fizzbuzz)
	defer close(fb.number)
	for i := 1; i <= fb.n; i++ {
		switch {
		case i%15 == 0:
			fb.fizzbuzz <- struct{}{}
			<-fb.number
		case i%5 == 0:
			fb.buzz <- struct{}{}
			<-fb.number
		case i%3 == 0:
			fb.fizz <- struct{}{}
			<-fb.number
		default:
			printNumber(i)
		}
	}
}
