package foo

import (
	"strings"
	"sync"
	"testing"
)

func TestFoo(t *testing.T) {
	tests := [][3]int{
		{2, 3, 1},
		{1, 3, 2},
		{3, 2, 1},
		{3, 1, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	f := NewFoo()
	for _, nums := range tests {
		output := f.run(nums)
		if output != "firstsecondthird" {
			t.Errorf("expected firstsecondthird, but got %s", output)
		}
	}
}

func (f *Foo) run(order [3]int) string {
	ch := make(chan string, 3)
	first := func() { ch <- "first" }
	second := func() { ch <- "second" }
	third := func() { ch <- "third" }

	var wg sync.WaitGroup
	wg.Add(3)
	for _, o := range order {
		switch o {
		case 1:
			go func() {
				defer wg.Done()
				f.First(first)
			}()
		case 2:
			go func() {
				defer wg.Done()
				f.Second(second)
			}()
		case 3:
			go func() {
				defer wg.Done()
				f.Third(third)
			}()
		}
	}
	wg.Wait()

	result := strings.Builder{}
	for i := 0; i < 3; i++ {
		result.WriteString(<-ch)
	}
	return result.String()
}
