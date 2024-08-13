package main

import (
	"fmt"
	"sync"
)

type FooBar struct {
	n   int
	foo chan struct{}
	bar chan struct{}
}

func NewFooBar(n int) *FooBar {
	fb := &FooBar{n: n}
	fb.foo = make(chan struct{})
	fb.bar = make(chan struct{})
	return fb
}

func (fb *FooBar) Foo(printFoo func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.foo
		printFoo()
		fb.bar <- struct{}{}
	}
}

func (fb *FooBar) Bar(printBar func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.bar
		printBar()
		if i < fb.n-1 {
			fb.foo <- struct{}{}
		}
	}
}

func (fb *FooBar) Run() {
	if fb.n <= 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		fb.Foo(func() { fmt.Print("foo") })
	}()
	go func() {
		defer wg.Done()
		fb.Bar(func() { fmt.Print("bar") })
	}()
	fb.foo <- struct{}{}
	wg.Wait()
}

func main() {
	fb := NewFooBar(3)
	fb.Run()
}
