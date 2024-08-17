package zeroevenodd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestZeroEvenOdd(t *testing.T) {
	ns := []int{2, 5, 10}
	for _, n := range ns {
		if err := validateN(n); err != nil {
			t.Fatalf("validateN(%d) = %v, want nil", n, err)
		}
		z := NewZeroEvenOdd(n)
		actual := z.run()
		expected := makeExpected(n)
		if actual != expected {
			t.Errorf("NewZeroEvenOdd(%d).run() = %q, want %q", n, actual, expected)
		}
	}
}

func validateN(n int) error {
	if n < 1 {
		return fmt.Errorf("n = %d, want n > 0", n)
	}
	if n > 1000 {
		return fmt.Errorf("n = %d, want n <= 1000", n)
	}
	return nil
}

func (z *ZeroEvenOdd) run() string {
	var wg sync.WaitGroup
	wg.Add(3)
	ch := make(chan int, 2*z.n)

	go func() {
		defer wg.Done()
		z.Zero(func(n int) {
			if n == 0 {
				ch <- n
			}
		})
	}()

	go func() {
		defer wg.Done()
		z.Even(func(n int) {
			if n&1 == 0 {
				ch <- n
			}
		})
	}()

	go func() {
		defer wg.Done()
		z.Odd(func(n int) {
			if n&1 == 1 {
				ch <- n
			}
		})
	}()

	wg.Wait()
	close(ch)

	var actual strings.Builder
	for n := range ch {
		actual.WriteString(strconv.Itoa(n))
	}
	return actual.String()
}

func makeExpected(n int) string {
	var expected strings.Builder
	for i := 1; i <= n; i++ {
		expected.WriteString("0")
		expected.WriteString(strconv.Itoa(i))
	}
	return expected.String()
}
