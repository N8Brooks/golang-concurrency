// Package testsuite contains reusable behavioral tests for multithreaded
// FizzBuzz implementations.
package testsuite

import (
	"slices"
	"strconv"
	"sync"
	"testing"
	"testing/synctest"
)

type FizzBuzz interface {
	Fizz(printFizz func())
	Buzz(printBuzz func())
	FizzBuzz(printFizzBuzz func())
	Number(printNumber func(int))
}

type worker struct {
	done chan struct{}
}

func startFizz(fb FizzBuzz, printFizz func()) worker {
	w := worker{done: make(chan struct{})}
	go func() {
		fb.Fizz(printFizz)
		close(w.done)
	}()
	return w
}

func startBuzz(fb FizzBuzz, printBuzz func()) worker {
	w := worker{done: make(chan struct{})}
	go func() {
		fb.Buzz(printBuzz)
		close(w.done)
	}()
	return w
}

func startFizzBuzz(fb FizzBuzz, printFizzBuzz func()) worker {
	w := worker{done: make(chan struct{})}
	go func() {
		fb.FizzBuzz(printFizzBuzz)
		close(w.done)
	}()
	return w
}

func startNumber(fb FizzBuzz, printNumber func(int)) worker {
	w := worker{done: make(chan struct{})}
	go func() {
		fb.Number(printNumber)
		close(w.done)
	}()
	return w
}

func requireClosed(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
	default:
		t.Fatal(message)
	}
}

func requireOpen[T any](t *testing.T, ch <-chan T, message string) {
	t.Helper()

	select {
	case <-ch:
		t.Fatal(message)
	default:
	}
}

func makeExpected(n int) []string {
	expected := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			expected = append(expected, "fizzbuzz")
		case i%5 == 0:
			expected = append(expected, "buzz")
		case i%3 == 0:
			expected = append(expected, "fizz")
		default:
			expected = append(expected, strconv.Itoa(i))
		}
	}
	return expected
}

func runSequence(t *testing.T, fb FizzBuzz, n int) []string {
	t.Helper()

	output := make(chan string, n)
	var mu sync.Mutex
	actual := make([]string, 0, n)

	workers := []worker{
		startFizz(fb, func() { output <- "fizz" }),
		startBuzz(fb, func() { output <- "buzz" }),
		startFizzBuzz(fb, func() { output <- "fizzbuzz" }),
		startNumber(fb, func(i int) { output <- strconv.Itoa(i) }),
	}

	synctest.Wait()

	for _, w := range workers {
		requireClosed(t, w.done, "worker did not complete")
	}

	close(output)
	for value := range output {
		mu.Lock()
		actual = append(actual, value)
		mu.Unlock()
	}

	return actual
}

func Run(t *testing.T, newImpl func(n int) FizzBuzz) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, n := range []int{1, 2, 3, 5, 15, 16, 30} {
			t.Run(strconv.Itoa(n), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					actual := runSequence(t, newImpl(n), n)
					expected := makeExpected(n)
					if !slices.Equal(actual, expected) {
						t.Fatalf("got %q, want %q", actual, expected)
					}
				})
			})
		}
	})

	t.Run("WorkersMayStartBeforeNumber", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 15
			fb := newImpl(n)
			output := make(chan string, n)

			fizz := startFizz(fb, func() { output <- "fizz" })
			buzz := startBuzz(fb, func() { output <- "buzz" })
			fizzbuzz := startFizzBuzz(fb, func() { output <- "fizzbuzz" })
			synctest.Wait()

			requireOpen(t, output, "workers produced output before the number worker started")
			requireOpen(t, fizz.done, "fizz worker exited before the number worker started")
			requireOpen(t, buzz.done, "buzz worker exited before the number worker started")
			requireOpen(t, fizzbuzz.done, "fizzbuzz worker exited before the number worker started")

			number := startNumber(fb, func(i int) { output <- strconv.Itoa(i) })
			synctest.Wait()

			requireClosed(t, fizz.done, "fizz worker did not complete")
			requireClosed(t, buzz.done, "buzz worker did not complete")
			requireClosed(t, fizzbuzz.done, "fizzbuzz worker did not complete")
			requireClosed(t, number.done, "number worker did not complete")

			close(output)
			actual := make([]string, 0, n)
			for value := range output {
				actual = append(actual, value)
			}
			expected := makeExpected(n)
			if !slices.Equal(actual, expected) {
				t.Fatalf("got %q, want %q", actual, expected)
			}
		})
	})
}

func Benchmark(b *testing.B, newImpl func(n int) FizzBuzz) {
	b.Helper()
	b.ReportAllocs()

	for _, n := range []int{15, 50} {
		b.Run("N"+strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				fb := newImpl(n)
				var wg sync.WaitGroup
				wg.Add(4)

				go func() {
					fb.Fizz(func() {})
					wg.Done()
				}()
				go func() {
					fb.Buzz(func() {})
					wg.Done()
				}()
				go func() {
					fb.FizzBuzz(func() {})
					wg.Done()
				}()
				go func() {
					fb.Number(func(int) {})
					wg.Done()
				}()

				wg.Wait()
			}
		})
	}
}
