package main

import "fmt"

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
		fb.foo <- struct{}{}
	}
}

func (fb *FooBar) Run() {
	fb.foo <- struct{}{}
}

func main() {
	fb := NewFooBar(10)
	go fb.Foo(func() { fmt.Print("foo") })
	go fb.Bar(func() { fmt.Print("bar") })
	fb.Run()
}
