//go:build challenge

package diningsavages_test

import (
	"testing"

	diningsavages "github.com/N8Brooks/golang-concurrency/dining_savages"
	"github.com/N8Brooks/golang-concurrency/dining_savages/testsuite"
)

func TestDiningSavages(t *testing.T) {
	testsuite.Run(t, func(capacity int) testsuite.DiningSavages {
		return diningsavages.NewDiningSavages(capacity)
	})
}

func BenchmarkDiningSavages(b *testing.B) {
	testsuite.Benchmark(b, func(capacity int) testsuite.DiningSavages {
		return diningsavages.NewDiningSavages(capacity)
	})
}
