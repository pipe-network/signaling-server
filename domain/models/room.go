package models

import "github.com/pipe-network/signaling-server/domain/dtos"

type Room struct {
	InitiatorsPublicKey dtos.Key
	clients             map[dtos.Address]*Client
}

func NewRoom(publicKey dtos.Key) *Room {
	return &Room{
		InitiatorsPublicKey: publicKey,
		clients:             map[dtos.Address]*Client{},
	}
}

// AddClient returns false if the client was already added, otherwise adds the client and returns true
func (r *Room) AddClient(client *Client) bool {
	if _, ok := r.clients[client.Address]; ok {
		return false
	}
	r.clients[client.Address] = client
	return true
}

// GetClient returns the client with given id
func (r *Room) GetClient(address dtos.Address) *Client {
	return r.clients[address]
}

// RemoveClient returns true if the client with given id was found and removed, otherwise false
func (r *Room) RemoveClient(address dtos.Address) bool {
	if _, ok := r.clients[address]; ok {
		delete(r.clients, address)
		return true
	}

	return false
}
