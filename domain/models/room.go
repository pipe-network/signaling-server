package models

import (
	"errors"
	"github.com/pipe-network/signaling-server/domain/values"
	"sync"
)

var (
	RoomFull           = errors.New("room full")
	InitiatiorNotFound = errors.New("initiator not found")
)

type Room struct {
	InitiatorsPublicKey          values.Key
	clients                      map[string]*Client
	reservedResponderAddresses   map[int]bool
	reserveResponderAddressMutex sync.Mutex
}

func NewRoom(publicKey values.Key) *Room {
	return &Room{
		InitiatorsPublicKey:          publicKey,
		clients:                      map[string]*Client{},
		reservedResponderAddresses:   initReservedResponderAddresses(),
		reserveResponderAddressMutex: sync.Mutex{},
	}
}

func initReservedResponderAddresses() map[int]bool {
	reservedResponderAddresses := map[int]bool{}
	for i := 2; i < int(values.MaxAddress); i++ {
		reservedResponderAddresses[i] = false
	}
	return reservedResponderAddresses
}

// AddClient returns false if the client was already added, otherwise adds the client and returns true
func (r *Room) AddClient(client *Client) bool {
	if _, ok := r.clients[client.ID]; ok {
		return false
	}
	r.clients[client.ID] = client
	return true
}

// RemoveClient returns true if the client with given id was found and removed, otherwise false
func (r *Room) RemoveClient(client *Client) bool {
	if _, ok := r.clients[client.ID]; ok {
		delete(r.clients, client.ID)
		return true
	}
	return false
}

func (r *Room) CountResponders() int {
	count := 0
	for _, client := range r.clients {
		if client.Address != values.InitiatorAddress && client.Address != values.UnassignedAddress {
			count++
		}
	}
	return count
}

func (r *Room) NextFreeResponderAddress() (*values.Address, error) {
	r.reserveResponderAddressMutex.Lock()
	defer r.reserveResponderAddressMutex.Unlock()
	for addressInt, reserved := range r.reservedResponderAddresses {
		if !reserved {
			address := values.Address(addressInt)
			return &address, nil
		}
	}
	return nil, RoomFull
}

func (r *Room) ReserveAddress(address values.Address) {
	r.reservedResponderAddresses[int(address)] = true
}

func (r *Room) ReleaseAddress(address values.Address) {
	r.reservedResponderAddresses[int(address)] = false
}

func (r *Room) KickCurrentInitiator() {
	for _, client := range r.clients {
		if client.Address == values.InitiatorAddress {
			client.DropConnection(values.DroppedByInitiatorCode)
			r.RemoveClient(client)
			return
		}
	}
}

func (r *Room) Responders() []*Client {
	var clients []*Client
	for _, client := range r.clients {
		if client.IsResponder() {
			clients = append(clients, client)
		}
	}
	return clients
}

func (r *Room) Client(address values.Address) *Client {
	for _, client := range r.clients {
		if client.Address == address {
			return client
		}
	}
	return nil
}

func (r *Room) Initiator() *Client {
	for _, client := range r.clients {
		if client.IsInitiator() {
			return client
		}
	}
	return nil
}
