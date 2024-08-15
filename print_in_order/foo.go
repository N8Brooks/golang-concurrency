package foo

import (
	"sync"
)

type Foo struct {
	firstJobDone  sync.Mutex
	secondJobDone sync.Mutex
}

func NewFoo() *Foo {
	f := &Foo{}
	f.firstJobDone.Lock()
	f.secondJobDone.Lock()
	return f
}

func (f *Foo) First(printFirst func()) {
	printFirst()
	f.firstJobDone.Unlock()
}

func (f *Foo) Second(printSecond func()) {
	f.firstJobDone.Lock()
	printSecond()
	f.secondJobDone.Unlock()
}

func (f *Foo) Third(printThird func()) {
	f.secondJobDone.Lock()
	printThird()
}
