package models

import (
	"crypto/rand"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/values"
	"golang.org/x/crypto/nacl/box"
	"sync"
)

type Client struct {
	Address           values.Address
	SessionPrivateKey values.Key
	SessionPublicKey  values.Key
	Cookie            values.Cookie
	SequenceNumber    values.SequenceNumber
	OverflowNumber    values.OverflowNumber

	state values.ClientState

	connection      *websocket.Conn
	connectionMutex *sync.Mutex
}

func NewClient(
	connection *websocket.Conn,
) (*Client, error) {
	sessionPublicKey, sessionPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	cookie, err := values.NewRandomCookie()
	if err != nil {
		return nil, err
	}

	sequenceNumber, err := values.NewSequenceNumber()
	if err != nil {
		return nil, err
	}

	return &Client{
		Address:           values.UnassignedAddress,
		SessionPrivateKey: *sessionPrivateKey,
		SessionPublicKey:  *sessionPublicKey,
		Cookie:            *cookie,
		SequenceNumber:    *sequenceNumber,
		OverflowNumber:    values.NewOverflowNumber(),
		connection:        connection,
		connectionMutex:   &sync.Mutex{},
	}, nil
}

func (c *Client) SendSignalingMessage(signalingMessage SignalingMessage) error {
	c.connectionMutex.Lock()
	defer c.connectionMutex.Unlock()
	signalingMessageBytes, err := signalingMessage.Bytes()
	if err != nil {
		return err
	}
	return c.connection.WriteMessage(websocket.BinaryMessage, signalingMessageBytes)
}
