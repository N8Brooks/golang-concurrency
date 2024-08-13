package main

import (
	"fmt"
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

func (f *Foo) Run(order [3]int) {
	var wg sync.WaitGroup
	wg.Add(3)

	for _, o := range order {
		switch o {
		case 1:
			go func() {
				defer wg.Done()
				f.First(func() { fmt.Print("first") })
			}()
		case 2:
			go func() {
				defer wg.Done()
				f.Second(func() { fmt.Print("second") })
			}()
		case 3:
			go func() {
				defer wg.Done()
				f.Third(func() { fmt.Print("third") })
			}()
		}
	}

	wg.Wait()
}

func main() {
	f := NewFoo()
	f.Run([3]int{2, 3, 1})
	fmt.Println()
	f.Run([3]int{1, 3, 2})
	fmt.Println()
	f.Run([3]int{3, 2, 1})
	fmt.Println()
	f.Run([3]int{3, 1, 2})
	fmt.Println()
	f.Run([3]int{1, 2, 3})
	fmt.Println()
	f.Run([3]int{2, 1, 3})
}
