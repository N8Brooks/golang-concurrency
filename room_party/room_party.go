//go:build challenge

// Package room_party contains the challenge version of the room party
// problem.
//
// Students may freely share a room, but the Dean of Students may only enter
// when the room is empty or when the crowd has reached the breakup threshold.
// While the Dean is in the room, no additional students may enter.
package room_party

import "context"

type RoomParty struct{}

func NewRoomParty() *RoomParty {
	return &RoomParty{}
}

func (r *RoomParty) Student(ctx context.Context, party func()) error {
	panic("unimplemented")
}

func (r *RoomParty) Dean(ctx context.Context, search, breakup func()) error {
	panic("unimplemented")
}
