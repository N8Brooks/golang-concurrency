//go:build challenge

// Package zeroevenodd contains the challenge version of the print zero even
// odd problem.
//
// Three goroutines cooperate to print the sequence `010203...0n`. The zero
// goroutine prints `0` before every number, the odd goroutine prints odd
// numbers in increasing order, and the even goroutine prints even numbers in
// increasing order.
package zeroevenodd

type ZeroEvenOdd struct {
	n int
}

func NewZeroEvenOdd(n int) *ZeroEvenOdd {
	return &ZeroEvenOdd{n: n}
}

func (z *ZeroEvenOdd) Zero(printNumber func(int)) {
	panic("unimplemented")
}

func (z *ZeroEvenOdd) Even(printNumber func(int)) {
	panic("unimplemented")
}

func (z *ZeroEvenOdd) Odd(printNumber func(int)) {
	panic("unimplemented")
}
