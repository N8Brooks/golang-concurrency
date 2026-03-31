// Package testsuite contains reusable behavioral tests for sleep implementations.
package testsuite

import (
	"testing"
	"time"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
)

type SleepFunc func(time.Duration) promise.Promiser[struct{}]

func Run(t *testing.T, sleep SleepFunc) {
	t.Helper()

	for _, d := range []time.Duration{50 * time.Millisecond, 200 * time.Millisecond} {
		t.Run(d.String(), func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			if _, err := sleep(d).Result(); err != nil {
				t.Fatalf("Sleep(%v) returned %v", d, err)
			}
			if elapsed := time.Since(start); elapsed < d {
				t.Fatalf("Sleep(%v) took %v", d, elapsed)
			}
		})
	}
}

func Benchmark(b *testing.B, sleep SleepFunc) {
	b.Helper()
	b.ReportAllocs()

	for b.Loop() {
		if _, err := sleep(0).Result(); err != nil {
			b.Fatalf("Sleep(0) returned %v", err)
		}
	}
}
