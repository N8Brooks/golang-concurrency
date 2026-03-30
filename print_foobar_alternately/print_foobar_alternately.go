//go:build challenge

// Package foobaralternately contains the challenge version of the print
// FooBar alternately problem.
//
// Given n, one thread should invoke printFoo n times and another should invoke
// printBar n times, producing the sequence "foobar" repeated n times.
package foobaralternately

type FooBar struct {
	n int
}

func NewFooBar(n int) *FooBar {
	return &FooBar{n: n}
}

func (fb *FooBar) Foo(printFoo func()) {
	panic("unimplemented")
}

func (fb *FooBar) Bar(printBar func()) {
	panic("unimplemented")
}
