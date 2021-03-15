package values

import (
	"golang.org/x/crypto/nacl/box"
)

type TypedMessage interface {
	MessageType() MessageType
}

type Message struct {
	Type MessageType `msgpack:"type"`
}

type ServerHelloMessage struct {
	Message
	Key Key `msgpack:"key"`
}

type ClientHelloMessage struct {
	Message
	Key Key `msgpack:"key"`
}

type ClientAuthMessage struct {
	Message
	YourCookie   Cookie   `msgpack:"your_cookie"`
	SubProtocols []string `msgpack:"subprotocols"`
	PingInterval int      `msgpack:"ping_interval"`
	YourKey      Key      `msgpack:"your_key"`
}

type ServerAuthMessage struct {
	Message
	YourCookie         Cookie    `msgpack:"your_cookie"`
	SignedKeys         []byte    `msgpack:"signed_keys"`
	InitiatorConnected bool      `msgpack:"initiator_connected"`
	Responders         []Address `msgpack:"responders"`
}

type NewInitiatorMessage struct {
	Message
}

type NewResponderMessage struct {
	Message
	ID Address `msgpack:"id"`
}

type DropResponderMessage struct {
	Message
	ID     Address   `msgpack:"id"`
	Reason CloseCode `msgpack:"reason"`
}

type DisconnectedMessage struct {
	Message
	ID Address `msgpack:"id"`
}

type SendErrorMessage struct {
	Message
	ID []byte `msgpack:"id"`
}

func (m *Message) MessageType() MessageType {
	return m.Type
}

func NewServerHelloMessage(sessionPublicKey Key) ServerHelloMessage {
	return ServerHelloMessage{
		Message: Message{
			Type: ServerHello,
		},
		Key: sessionPublicKey,
	}
}

func NewServerAuthMessage(
	incomingCookie Cookie,
	outgoingSessionPublicKey Key,
	clientPermanentPublicKey Key,
	serverPermanentPrivateKey Key,
	nonce Nonce,
	initiatorConnected bool,
	responderAddresses []Address,
) ServerAuthMessage {
	var signedKeys []byte
	nonceBytes := nonce.Bytes()
	peersPublicKeyBytes := clientPermanentPublicKey.Bytes()
	privateKeyBytes := serverPermanentPrivateKey.Bytes()
	signedKeys = append(signedKeys, outgoingSessionPublicKey[:]...)
	signedKeys = append(signedKeys, clientPermanentPublicKey[:]...)
	signedKeys = box.Seal(nil, signedKeys, &nonceBytes, &peersPublicKeyBytes, &privateKeyBytes)

	return ServerAuthMessage{
		Message: Message{
			Type: ServerAuth,
		},
		YourCookie:         incomingCookie,
		SignedKeys:         signedKeys,
		InitiatorConnected: initiatorConnected,
		Responders:         responderAddresses,
	}
}

func NewNewInitiatorMessage() NewInitiatorMessage {
	return NewInitiatorMessage{
		Message: Message{Type: NewInitiator},
	}
}

func NewNewResponderMessage(responderAddress Address) NewResponderMessage {
	return NewResponderMessage{
		Message: Message{Type: NewResponder},
		ID:      responderAddress,
	}
}

func (m ClientAuthMessage) ContainsSubProtocol(subProtocol string) bool {
	for _, tempSubProtocol := range m.SubProtocols {
		if subProtocol == tempSubProtocol {
			return true
		}
	}
	return false
}
