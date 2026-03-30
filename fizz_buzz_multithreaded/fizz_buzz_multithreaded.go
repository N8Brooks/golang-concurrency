//go:build challenge

// Package fizzbuzz contains the challenge version of the multithreaded FizzBuzz
// problem.
//
// Four threads cooperate to print the FizzBuzz sequence from 1 through n. Each
// worker may only emit its own category, and together they must produce the
// exact ordered sequence.
package fizzbuzz

type FizzBuzz struct {
	n int
}

func NewFizzBuzz(n int) *FizzBuzz {
	return &FizzBuzz{n: n}
}

func (fb *FizzBuzz) Fizz(printFizz func()) {
	panic("unimplemented")
}

func (fb *FizzBuzz) Buzz(printBuzz func()) {
	panic("unimplemented")
}

func (fb *FizzBuzz) FizzBuzz(printFizzBuzz func()) {
	panic("unimplemented")
}

func (fb *FizzBuzz) Number(printNumber func(int)) {
	panic("unimplemented")
}
