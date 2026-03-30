package solutions

type Channel struct {
	n   int
	foo chan struct{}
	bar chan struct{}
}

func NewChannel(n int) *Channel {
	foo := make(chan struct{}, 1)
	foo <- struct{}{}
	return &Channel{
		n:   n,
		foo: foo,
		bar: make(chan struct{}, 1),
	}
}

func (fb *Channel) Foo(printFoo func()) {
	for range fb.n {
		<-fb.foo
		printFoo()
		fb.bar <- struct{}{}
	}
}

func (fb *Channel) Bar(printBar func()) {
	for range fb.n {
		<-fb.bar
		printBar()
		fb.foo <- struct{}{}
	}
}
