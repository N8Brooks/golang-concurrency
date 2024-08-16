package fizzbuzz

import (
	"errors"
	"slices"
	"strconv"
	"sync"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	ns := []int{
		15,
		5,
	}
	for _, n := range ns {
		fb := NewFizzBuzz(n)
		actual, err := fb.run()
		if err != nil {
			t.Fatalf("NewFizzBuzz(%d).run() = %v, want nil", n, err)
		}
		expected := makeExpected(n)
		if !slices.Equal(actual, expected) {
			t.Errorf("NewFizzBuzz(%d).run() = %q; want %q", n, actual, expected)
		}
	}
}

func (fb *FizzBuzz) run() ([]string, error) {
	if err := validateN(fb.n); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(4)

	ch := make(chan string, fb.n)

	go func() {
		defer wg.Done()
		fb.Fizz(func() { ch <- "fizz" })
	}()

	go func() {
		defer wg.Done()
		fb.Buzz(func() { ch <- "buzz" })
	}()

	go func() {
		defer wg.Done()
		fb.FizzBuzz(func() { ch <- "fizzbuzz" })
	}()

	go func() {
		defer wg.Done()
		fb.Number(func(i int) { ch <- strconv.Itoa(i) })
	}()

	wg.Wait()
	close(ch)

	actual := make([]string, 0, fb.n)
	for i := 0; i < fb.n; i++ {
		actual = append(actual, <-ch)
	}
	return actual, nil
}

func validateN(n int) error {
	if n < 1 {
		return errors.New("n must be greater than 0")
	}
	if n > 50 {
		return errors.New("n must be less than or equal to 50")
	}
	return nil
}

func makeExpected(n int) []string {
	var expected []string
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			expected = append(expected, "fizzbuzz")
		case i%3 == 0:
			expected = append(expected, "fizz")
		case i%5 == 0:
			expected = append(expected, "buzz")
		default:
			expected = append(expected, strconv.Itoa(i))
		}
	}
	return expected
}
