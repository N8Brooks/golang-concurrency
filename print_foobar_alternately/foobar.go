package foobar

type FooBar struct {
	n   int
	foo chan struct{}
	bar chan struct{}
}

func NewFooBar(n int) *FooBar {
	fb := &FooBar{n: n}
	fb.foo = make(chan struct{}, 1)
	fb.bar = make(chan struct{})
	fb.foo <- struct{}{}
	return fb
}

func (fb *FooBar) Foo(printFoo func()) {
	defer close(fb.bar)
	for _ = range fb.foo {
		printFoo()
		fb.bar <- struct{}{}
	}
}

func (fb *FooBar) Bar(printBar func()) {
	defer close(fb.foo)
	for i := 1; i < fb.n; i++ {
		<-fb.bar
		printBar()
		fb.foo <- struct{}{}
	}
	<-fb.bar
	printBar()
}
