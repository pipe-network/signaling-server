package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/values"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/box"
	"sync"
	"time"
)

var (
	DefaultPongWait, _ = time.ParseDuration("30s")
	NotAllowedToRelay  = func(destinationAddress values.Address) error {
		return errors.New(fmt.Sprintf("not allowed to relay messages to %x", destinationAddress))
	}
	IdentitiesDoNotMatch = func(clientAddress, sourceAddress values.Address) error {
		return errors.New(
			fmt.Sprintf("identities do not match, expected %x, got %x", clientAddress, sourceAddress),
		)
	}
	InvalidCookie         = errors.New("invalid cookie")
	InvalidSequenceNumber = errors.New("invalid sequence number")
)

type Client struct {
	ID      string
	Address values.Address

	SessionPrivateKey values.Key
	SessionPublicKey  values.Key

	PermanentPublicKey values.Key

	OutgoingCookie         values.Cookie
	OutgoingSequenceNumber values.SequenceNumber
	OutgoingOverflowNumber values.OverflowNumber

	IncomingCookie         values.Cookie
	IncomingSequenceNumber values.SequenceNumber
	IncomingOverflowNumber values.OverflowNumber

	state      values.ClientState
	pingTicker *time.Ticker

	connection           *websocket.Conn
	connected            bool
	connectionWriteMutex *sync.Mutex
}

func NewClient(
	connection *websocket.Conn,
	room *Room,
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
		ID:                     uuid.NewString(),
		Address:                values.UnassignedAddress,
		SessionPrivateKey:      *sessionPrivateKey,
		SessionPublicKey:       *sessionPublicKey,
		OutgoingCookie:         *cookie,
		OutgoingSequenceNumber: *sequenceNumber,
		OutgoingOverflowNumber: values.NewOverflowNumber(),
		connection:             connection,
		connected:              true,
		connectionWriteMutex:   &sync.Mutex{},
	}, nil
}

func (c *Client) Nonce() values.Nonce {
	return values.Nonce{
		Cookie:         c.OutgoingCookie,
		Source:         values.ServerAddress,
		Destination:    c.Address,
		SequenceNumber: c.OutgoingSequenceNumber,
		OverflowNumber: c.OutgoingOverflowNumber,
	}
}

func (c *Client) OutgoingCombinedSequenceNumber() values.CombinedSequenceNumber {
	return values.NewCombinedSequenceNumber(c.OutgoingSequenceNumber, c.OutgoingOverflowNumber)
}

func (c *Client) IncomingCombinedSequenceNumber() values.CombinedSequenceNumber {
	return values.NewCombinedSequenceNumber(c.IncomingSequenceNumber, c.IncomingOverflowNumber)
}

func (c *Client) IsP2PAllowed(destinationAddress values.Address) bool {
	return c.IsAuthenticated() && c.Address != destinationAddress
}

func (c *Client) IncomingNonceEmpty() bool {
	return c.IncomingCookie.Empty() && c.IncomingOverflowNumber.Empty() && c.IncomingSequenceNumber.Empty()
}

func (c *Client) SetIncomingNonce(nonce values.Nonce) {
	c.IncomingCookie = nonce.Cookie
	c.IncomingSequenceNumber = nonce.SequenceNumber
	c.IncomingOverflowNumber = nonce.OverflowNumber
}

func (c *Client) SetPermanentPublicKey(permanentPublicKey values.Key) {
	c.PermanentPublicKey = permanentPublicKey
}

func (c *Client) SetAddress(address values.Address) {
	c.Address = address
}

func (c *Client) ValidateNonce(nonce values.Nonce) error {
	isAddressedToServer := nonce.Destination == values.ServerAddress
	combinedSequenceNumber := values.NewCombinedSequenceNumber(nonce.SequenceNumber, nonce.OverflowNumber)

	// Validate destination address
	if !isAddressedToServer && !c.IsP2PAllowed(nonce.Destination) {
		return NotAllowedToRelay(nonce.Destination)
	}

	// Validate source address
	if nonce.Source != c.Address {
		return IdentitiesDoNotMatch(c.Address, nonce.Source)
	}

	if isAddressedToServer {
		if !c.IsCookieValid(nonce.Cookie) {
			return InvalidCookie
		}

		if !c.IsCombinedSequenceNumberValid(combinedSequenceNumber) {
			return InvalidSequenceNumber
		}
	}
	return nil
}

func (c *Client) IsCookieValid(cookie values.Cookie) bool {
	if c.IncomingCookie.Empty() {
		// Check that the outgoing cookie is not the same as the new cookie received from client
		if c.OutgoingCookie.Equal(cookie) {
			return false
		}
	} else {
		// Check that the cookie hasn't changed
		if !c.IncomingCookie.Equal(cookie) {
			return false
		}
	}
	return true
}

func (c *Client) AssignToInitiator() {
	c.SetAddress(values.InitiatorAddress)
}

func (c *Client) MarkAsAuthenticated() {
	c.state = values.Authenticated
}

func (c *Client) IsCombinedSequenceNumberValid(combinedSequenceNumber values.CombinedSequenceNumber,
) bool {
	if c.IncomingCombinedSequenceNumber().Empty() && combinedSequenceNumber.OverflowNumber.Int() != 0 {
		return false
	}

	if !combinedSequenceNumber.Equal(c.IncomingCombinedSequenceNumber()) {
		return false
	}

	return true
}

func (c *Client) DropConnection(code values.CloseCode) {
	c.connectionWriteMutex.Lock()
	c.connected = false
	defer c.connectionWriteMutex.Unlock()
	_ = c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code.Int(), code.Message()))
	_ = c.CloseConnection()
}

func (c *Client) CloseConnection() error {
	return c.connection.Close()
}

func (c *Client) ReadMessage() (messageType int, p []byte, err error) {
	if !c.connected {
		return 0, []byte{}, errors.New("client not connected anymore")
	}
	return c.connection.ReadMessage()
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
	outgoingCombinedSequenceNumber := c.OutgoingCombinedSequenceNumber()
	incrementedOutgoingCombinedSequenceNumber, err := outgoingCombinedSequenceNumber.Increment()
	if err != nil {
		return err
	}
	c.OutgoingSequenceNumber = incrementedOutgoingCombinedSequenceNumber.SequenceNumber
	c.OutgoingOverflowNumber = incrementedOutgoingCombinedSequenceNumber.OverflowNumber
	return nil
}

func (c *Client) IsAuthenticated() bool {
	return c.state == values.Authenticated
}

func (c *Client) IsInitiator() bool {
	return c.Address == values.InitiatorAddress
}

func (c *Client) IsResponder() bool {
	return c.Address != values.InitiatorAddress && c.Address != values.UnassignedAddress
}

func (c *Client) SendBytes(bytes []byte) error {
	c.connectionWriteMutex.Lock()
	defer c.connectionWriteMutex.Unlock()
	err := c.connection.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		return err
	}
	err = c.IncrementOutgoingCombinedSequenceNumber()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PingTicker(pingPeriod time.Duration, pongWait time.Duration) {
	var err error
	c.pingTicker = time.NewTicker(pingPeriod)
	c.connection.SetPongHandler(
		func(string) error {
			log.Infof("Receiving pong from: %d", c.Address)
			return nil
		},
	)
	defer c.pingTicker.Stop()
	for {
		select {
		case <-c.pingTicker.C:
			log.Infof("Sending ping to: %d", c.Address)
			c.connectionWriteMutex.Lock()
			err = c.connection.SetWriteDeadline(time.Now().Add(pingPeriod))
			if err != nil {
				return
			}
			err = c.connection.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
			c.connectionWriteMutex.Unlock()
		}
	}
}

func (c *Client) Flush() {
	c.pingTicker.Stop()
}
