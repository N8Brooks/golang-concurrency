//go:build challenge

package debounce_test

import (
	"testing"

	debouncepkg "github.com/N8Brooks/golang-concurrency/debounce"
	"github.com/N8Brooks/golang-concurrency/debounce/testsuite"
)

func TestDebounce(t *testing.T) {
	testsuite.Run(t, debouncepkg.Debounce)
}

func BenchmarkDebounce(b *testing.B) {
	testsuite.Benchmark(b, debouncepkg.Debounce)
}
