package models

import (
	"bytes"
	"crypto/rand"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/values"
	"golang.org/x/crypto/nacl/box"
	"sync"
)

type Client struct {
	Address                values.Address
	SessionPrivateKey      values.Key
	SessionPublicKey       values.Key
	OutgoingCookie         values.Cookie
	OutgoingSequenceNumber values.SequenceNumber
	OutgoingOverflowNumber values.OverflowNumber

	IncomingCookie         values.Cookie
	IncomingSequenceNumber values.SequenceNumber
	IncomingOverflowNumber values.OverflowNumber

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
		Address:                values.UnassignedAddress,
		SessionPrivateKey:      *sessionPrivateKey,
		SessionPublicKey:       *sessionPublicKey,
		OutgoingCookie:         *cookie,
		OutgoingSequenceNumber: *sequenceNumber,
		OutgoingOverflowNumber: values.NewOverflowNumber(),
		connection:             connection,
		connectionMutex:        &sync.Mutex{},
	}, nil
}

func (c *Client) OutgoingCombinedSequenceNumber() values.CombinedSequenceNumber {
	return values.NewCombinedSequenceNumber(c.OutgoingSequenceNumber, c.OutgoingOverflowNumber)
}

func (c *Client) IncomingCombinedSequenceNumber() values.CombinedSequenceNumber {
	return values.NewCombinedSequenceNumber(c.IncomingSequenceNumber, c.IncomingOverflowNumber)
}

func (c *Client) IsP2PAllowed(destinationAddress values.Address) bool {
	return c.state == values.Authenticated && c.Address != destinationAddress
}

func (c *Client) IsCookieValid(cookie values.Cookie) bool {
	emptyCookie := values.Cookie{}
	if bytes.Equal(c.IncomingCookie[:], emptyCookie[:]) {
		if bytes.Equal(c.OutgoingCookie[:], cookie[:]) {
			return false
		}

		c.IncomingCookie = cookie
	} else {
		if !bytes.Equal(c.IncomingCookie[:], cookie[:]) {
			return false
		}
	}
	return true
}

func (c *Client) IsCombinedSequenceNumberValid(
	sequenceNumber values.SequenceNumber,
	overflowNumber values.OverflowNumber,
) bool {
	combinedSequenceNumber := values.NewCombinedSequenceNumber(sequenceNumber, overflowNumber)

	if c.IncomingCombinedSequenceNumber().Empty() {
		if overflowNumber.Int() != 0 {
			return false
		}
		c.IncomingSequenceNumber = sequenceNumber
		c.IncomingOverflowNumber = overflowNumber
	}

	if !combinedSequenceNumber.Equal(c.IncomingCombinedSequenceNumber()) {
		return false
	}

	return true
}

func (c *Client) IncrementIncomingCombinedSequenceNumber() error {
	incomingCombinedSequenceNumber := c.IncomingCombinedSequenceNumber()
	incrementedIncomingCombinedSequenceNumber, err := incomingCombinedSequenceNumber.Increment()
	if err != nil {
		return err
	}
	c.IncomingSequenceNumber = incrementedIncomingCombinedSequenceNumber.SequenceNumber
	c.IncomingOverflowNumber = incrementedIncomingCombinedSequenceNumber.OverflowNumber
	return nil
}

func (c *Client) IncrementOutgoingCombinedSequenceNumber() error {
	outgoingCombinedSequenceNumber := c.IncomingCombinedSequenceNumber()
	incrementedOutgoingCombinedSequenceNumber, err := outgoingCombinedSequenceNumber.Increment()
	if err != nil {
		return err
	}
	c.IncomingSequenceNumber = incrementedOutgoingCombinedSequenceNumber.SequenceNumber
	c.IncomingOverflowNumber = incrementedOutgoingCombinedSequenceNumber.OverflowNumber
	return nil
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
