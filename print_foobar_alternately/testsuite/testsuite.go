// Package testsuite contains reusable behavioral tests for FooBar implementations.
package testsuite

import (
	"strconv"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
)

type FooBar interface {
	Foo(printFoo func())
	Bar(printBar func())
}

func runOneShot(t *testing.T, fb FooBar, n int, fooFirst bool) {
	t.Helper()

	output := make(chan string, 2*n)
	doneFoo := make(chan struct{})
	doneBar := make(chan struct{})

	runFoo := func() {
		fb.Foo(func() {
			output <- "foo"
		})
		close(doneFoo)
	}
	runBar := func() {
		fb.Bar(func() {
			output <- "bar"
		})
		close(doneBar)
	}

	if fooFirst {
		go runFoo()
		synctest.Wait()
		go runBar()
	} else {
		go runBar()
		synctest.Wait()
		go runFoo()
	}

	synctest.Wait()

	select {
	case <-doneFoo:
	default:
		t.Fatal("Foo did not complete")
	}
	select {
	case <-doneBar:
	default:
		t.Fatal("Bar did not complete")
	}

	var builder strings.Builder
	for i := 0; i < 2*n; i++ {
		builder.WriteString(<-output)
	}

	expected := strings.Repeat("foobar", n)
	if actual := builder.String(); actual != expected {
		t.Fatalf("got %q, want %q", actual, expected)
	}
}

func Run(t *testing.T, newImpl func(n int) FooBar) {
	t.Helper()

	for _, tc := range []struct {
		name     string
		fooFirst bool
	}{
		{name: "FooFirst", fooFirst: true},
		{name: "BarFirst", fooFirst: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, n := range []int{1, 2, 4, 8} {
				t.Run("N"+strconv.Itoa(n), func(t *testing.T) {
					synctest.Test(t, func(t *testing.T) {
						runOneShot(t, newImpl(n), n, tc.fooFirst)
					})
				})
			}
		})
	}
}

func Benchmark(b *testing.B, newImpl func(n int) FooBar) {
	b.Helper()
	b.ReportAllocs()

	for _, n := range []int{1, 8, 64} {
		b.Run("N"+strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				fb := newImpl(n)
				var wg sync.WaitGroup
				wg.Add(2)

				go func() {
					defer wg.Done()
					fb.Foo(func() {})
				}()

				go func() {
					defer wg.Done()
					fb.Bar(func() {})
				}()

				wg.Wait()
			}
		})
	}
}
