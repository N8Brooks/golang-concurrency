// Package testsuite contains reusable behavioral tests for print zero even odd implementations.
package testsuite

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type ZeroEvenOdd interface {
	Zero(printNumber func(int))
	Even(printNumber func(int))
	Odd(printNumber func(int))
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

func runOne(t *testing.T, z ZeroEvenOdd, n int) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(3)
	ch := make(chan int, 2*n)

	go func() {
		defer wg.Done()
		z.Zero(func(v int) {
			if v != 0 {
				t.Errorf("Zero printed %d, want 0", v)
			}
			ch <- v
		})
	}()

	go func() {
		defer wg.Done()
		z.Even(func(v int) {
			if v == 0 || v&1 != 0 {
				t.Errorf("Even printed %d, want a positive even number", v)
			}
			ch <- v
		})
	}()

	go func() {
		defer wg.Done()
		z.Odd(func(v int) {
			if v&1 != 1 {
				t.Errorf("Odd printed %d, want an odd number", v)
			}
			ch <- v
		})
	}()

	wg.Wait()
	close(ch)

	var actual strings.Builder
	actual.Grow(2 * n)
	for v := range ch {
		actual.WriteString(strconv.Itoa(v))
	}

	expected := makeExpected(n)
	if actual.String() != expected {
		t.Fatalf("got %q, want %q", actual.String(), expected)
	}
}

func makeExpected(n int) string {
	var expected strings.Builder
	for i := 1; i <= n; i++ {
		expected.WriteByte('0')
		expected.WriteString(strconv.Itoa(i))
	}
	return expected.String()
}

func Run(t *testing.T, newImpl func(n int) ZeroEvenOdd) {
	t.Helper()

	for _, n := range []int{1, 2, 5, 10} {
		t.Run(fmt.Sprintf("N%d", n), func(t *testing.T) {
			if err := validateN(n); err != nil {
				t.Fatalf("validateN(%d) = %v, want nil", n, err)
			}
			runOne(t, newImpl(n), n)
		})
	}
}

func Benchmark(b *testing.B, newImpl func(n int) ZeroEvenOdd) {
	b.Helper()
	b.ReportAllocs()

	for _, n := range []int{8, 64} {
		b.Run(fmt.Sprintf("N%d", n), func(b *testing.B) {
			for b.Loop() {
				z := newImpl(n)
				var wg sync.WaitGroup
				wg.Add(3)

				go func() {
					defer wg.Done()
					z.Zero(func(int) {})
				}()

				go func() {
					defer wg.Done()
					z.Even(func(int) {})
				}()

				go func() {
					defer wg.Done()
					z.Odd(func(int) {})
				}()

				wg.Wait()
			}
		})
	}
}
