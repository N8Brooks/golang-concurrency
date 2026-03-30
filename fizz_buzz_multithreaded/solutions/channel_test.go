package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/fizz_buzz_multithreaded/solutions"
	"github.com/N8Brooks/golang-concurrency/fizz_buzz_multithreaded/testsuite"
)

func TestChannel(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.FizzBuzz {
		return solutions.NewChannel(n)
	})
}

func BenchmarkChannel(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.FizzBuzz {
		return solutions.NewChannel(n)
	})
}
