package main

import (
	"fmt"
	"sync"
)

type ZeroEvenOdd struct {
	n    int
	zero chan struct{}
	even chan struct{}
	odd  chan struct{}
}

func NewZeroEvenOdd(n int) *ZeroEvenOdd {
	z := &ZeroEvenOdd{n: n}
	z.zero = make(chan struct{})
	z.even = make(chan struct{})
	z.odd = make(chan struct{})
	return z
}

func (z *ZeroEvenOdd) Zero(printNumber func(int)) {
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

func (z *ZeroEvenOdd) Even(printNumber func(int)) {
	for i := 2; i <= z.n; i += 2 {
		<-z.even
		printNumber(i)
		if i < z.n {
			z.zero <- struct{}{}
		}
	}
}

func (z *ZeroEvenOdd) Odd(printNumber func(int)) {
	for i := 1; i <= z.n; i += 2 {
		<-z.odd
		printNumber(i)
		if i < z.n {
			z.zero <- struct{}{}
		}
	}
}

func (z *ZeroEvenOdd) Run() {
	if z.n <= 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		z.Zero(func(n int) { fmt.Print(n) })
	}()
	go func() {
		defer wg.Done()
		z.Even(func(n int) { fmt.Print(n) })
	}()
	go func() {
		defer wg.Done()
		z.Odd(func(n int) { fmt.Print(n) })
	}()
	z.zero <- struct{}{}
	wg.Wait()
}

func main() {
	z := NewZeroEvenOdd(2)
	z.Run()
}
