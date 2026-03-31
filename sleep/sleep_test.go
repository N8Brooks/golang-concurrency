//go:build challenge

package sleep_test

import (
	"testing"

	sleeppkg "github.com/N8Brooks/golang-concurrency/sleep"
	"github.com/N8Brooks/golang-concurrency/sleep/testsuite"
)

func TestSleep(t *testing.T) {
	testsuite.Run(t, sleeppkg.Sleep)
}

func BenchmarkSleep(b *testing.B) {
	testsuite.Benchmark(b, sleeppkg.Sleep)
}
