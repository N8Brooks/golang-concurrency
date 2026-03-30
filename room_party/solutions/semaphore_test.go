package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/room_party/solutions"
	"github.com/N8Brooks/golang-concurrency/room_party/testsuite"
)

func TestSemaphore(t *testing.T) {
	testsuite.Run(t, func() testsuite.RoomParty {
		return solutions.NewSemaphore()
	})
}

func BenchmarkSemaphore(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.RoomParty {
		return solutions.NewSemaphore()
	})
}
