package models

import "github.com/pipe-network/signaling-server/domain/dtos"

type Rooms struct {
	rooms map[dtos.Key]*Room
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: map[dtos.Key]*Room{},
	}
}

func (r *Rooms) AddRoom(room *Room) bool {
	if _, ok := r.rooms[room.InitiatorsPublicKey]; ok {
		return false
	}

	r.rooms[room.InitiatorsPublicKey] = room
	return true
}

func (r *Rooms) GetRoom(initiatorsPublicKey dtos.Key) *Room {
	return r.rooms[initiatorsPublicKey]
}

func (r *Rooms) RemoveRoom(initiatorsPublicKey dtos.Key) bool {
	if _, ok := r.rooms[initiatorsPublicKey]; ok {
		delete(r.rooms, initiatorsPublicKey)
		return true
	}
	return false
}
