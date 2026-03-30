//go:build challenge

package room_party_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/room_party"
	"github.com/N8Brooks/golang-concurrency/room_party/testsuite"
)

func TestRoomParty(t *testing.T) {
	testsuite.Run(t, func() testsuite.RoomParty {
		return room_party.NewRoomParty()
	})
}

func BenchmarkRoomParty(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.RoomParty {
		return room_party.NewRoomParty()
	})
}
