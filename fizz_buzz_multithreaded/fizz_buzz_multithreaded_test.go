//go:build challenge

package fizzbuzz_test

import (
	"testing"

	fizzbuzz "github.com/N8Brooks/golang-concurrency/fizz_buzz_multithreaded"
	"github.com/N8Brooks/golang-concurrency/fizz_buzz_multithreaded/testsuite"
)

func TestFizzBuzz(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.FizzBuzz {
		return fizzbuzz.NewFizzBuzz(n)
	})
}

func BenchmarkFizzBuzz(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.FizzBuzz {
		return fizzbuzz.NewFizzBuzz(n)
	})
}
