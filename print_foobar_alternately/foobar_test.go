package foobar

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestFooBar(t *testing.T) {
	ns := []int{
		1,
		2,
		4,
	}
	for _, n := range ns {
		if err := validateN(n); err != nil {
			t.Errorf("validateN(%d) = %q; want nil", n, err)
		}
		fb := NewFooBar(n)
		actual := fb.run()
		expected := strings.Repeat("foobar", n)
		if actual != expected {
			t.Errorf("NewFooBar(%d).run() = %q; want %q", n, actual, expected)
		}
	}
}

func validateN(n int) error {
	if n < 1 {
		return errors.New("n must be greater than 0")
	}
	if n > 1000 {
		return errors.New("n must be less than or equal to 1000")
	}
	return nil
}

func (fb *FooBar) run() string {
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan string, 2*fb.n)
	defer close(ch)

	go func() {
		defer wg.Done()
		fb.Foo(func() { ch <- "foo" })
	}()

	go func() {
		defer wg.Done()
		fb.Bar(func() { ch <- "bar" })
	}()

	wg.Wait()

	var sb strings.Builder
	for i := 0; i < 2*fb.n; i++ {
		sb.WriteString(<-ch)
	}
	return sb.String()
}
