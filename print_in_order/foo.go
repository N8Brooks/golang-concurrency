package foo

type Foo struct {
	first  chan struct{}
	second chan struct{}
}

func NewFoo() *Foo {
	f := &Foo{
		first:  make(chan struct{}),
		second: make(chan struct{}),
	}
	return f
}

func (f *Foo) First(printFirst func()) {
	printFirst()
	f.first <- struct{}{}
}

func (f *Foo) Second(printSecond func()) {
	<-f.first
	printSecond()
	f.second <- struct{}{}
}

func (f *Foo) Third(printThird func()) {
	<-f.second
	printThird()
}
