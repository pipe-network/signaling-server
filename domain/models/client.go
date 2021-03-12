package models

import (
	"crypto/rand"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/dtos"
	"golang.org/x/crypto/nacl/box"
	"sync"
)

type Client struct {
	Address           dtos.Address
	SessionPrivateKey dtos.Key
	SessionPublicKey  dtos.Key
	Cookie            dtos.Cookie

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

	cookie, err := dtos.NewRandomCookie()
	if err != nil {
		return nil, err
	}

	return &Client{
		SessionPrivateKey: *sessionPrivateKey,
		SessionPublicKey:  *sessionPublicKey,
		Cookie:            *cookie,
		Address:           dtos.Unassigned,
		connection:        connection,
		connectionMutex:   &sync.Mutex{},
	}, nil
}

func (c *Client) SendSignalingMessage(signalingMessage dtos.SignalingMessage) error {
	c.connectionMutex.Lock()
	defer c.connectionMutex.Unlock()
	signalingMessageBytes, err := signalingMessage.Bytes()
	if err != nil {
		return err
	}
	return c.connection.WriteMessage(websocket.BinaryMessage, signalingMessageBytes)
}
