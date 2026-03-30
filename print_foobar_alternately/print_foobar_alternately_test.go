//go:build challenge

package foobaralternately_test

import (
	"testing"

	foobaralternately "github.com/N8Brooks/golang-concurrency/print_foobar_alternately"
	"github.com/N8Brooks/golang-concurrency/print_foobar_alternately/testsuite"
)

func TestFooBar(t *testing.T) {
	testsuite.Run(t, func(n int) testsuite.FooBar {
		return foobaralternately.NewFooBar(n)
	})
}

func BenchmarkFooBar(b *testing.B) {
	testsuite.Benchmark(b, func(n int) testsuite.FooBar {
		return foobaralternately.NewFooBar(n)
	})
}
