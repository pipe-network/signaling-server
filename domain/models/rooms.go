package models

import (
	"github.com/pipe-network/signaling-server/domain/values"
)

type Rooms struct {
	rooms map[values.Key]*Room
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: map[values.Key]*Room{},
	}
}

func (r *Rooms) Size() int {
	return len(r.rooms)
}

func (r *Rooms) AddRoom(room *Room) bool {
	if _, ok := r.rooms[room.InitiatorsPublicKey]; ok {
		return false
	}

	r.rooms[room.InitiatorsPublicKey] = room
	return true
}

func (r *Rooms) GetRoom(initiatorsPublicKey values.Key) *Room {
	return r.rooms[initiatorsPublicKey]
}

func (r *Rooms) RemoveRoom(initiatorsPublicKey values.Key) bool {
	if _, ok := r.rooms[initiatorsPublicKey]; ok {
		delete(r.rooms, initiatorsPublicKey)
		return true
	}
	return false
}

func (r *Rooms) GetOrCreateRoom(initiatorsPublicKey values.Key) *Room {
	room := r.GetRoom(initiatorsPublicKey)
	if room == nil {
		room = NewRoom(initiatorsPublicKey)
		r.AddRoom(room)
	}
	return room
}
