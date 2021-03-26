package values

import (
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/nacl/box"
)

var (
	DecryptionFailed = errors.New("decryption failed")
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
	InitiatorConnected *bool      `msgpack:"initiator_connected,omitempty"`
	Responders         *[]Address `msgpack:"responders"`
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

func DecryptMessage(
	data []byte,
	nonce [NonceByteLength]byte,
	publicKey,
	privateKey Key,
) ([]byte, error) {
	publicKeyBytes := publicKey.Bytes()
	privateKeyBytes := privateKey.Bytes()

	decodedDataBytes, ok := box.Open(nil, data, &nonce, &publicKeyBytes, &privateKeyBytes)
	if !ok {
		return nil, DecryptionFailed
	}
	return decodedDataBytes, nil
}

func DecodeClientAuthMessageFromBytes(
	data []byte,
	nonce [NonceByteLength]byte,
	publicKey,
	privateKey Key,
) (*ClientAuthMessage, error) {
	var clientAuthMessage ClientAuthMessage
	var err error
	decryptedMessage, err := DecryptMessage(data, nonce, publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	err = msgpack.Unmarshal(decryptedMessage, &clientAuthMessage)
	if err != nil {
		return nil, err
	}
	return &clientAuthMessage, nil
}

func DecodeDropResponderMessageFromBytes(
	data []byte,
	nonce [NonceByteLength]byte,
	publicKey,
	privateKey Key,
) (*DropResponderMessage, error) {
	var dropResponderMessage DropResponderMessage
	decryptedMessage, err := DecryptMessage(data, nonce, publicKey, privateKey)
	err = msgpack.Unmarshal(decryptedMessage, &dropResponderMessage)
	if err != nil {
		return nil, err
	}
	return &dropResponderMessage, nil
}

func DecodeClientHelloMessageFromBytes(rawBytes []byte) (*ClientHelloMessage, error) {
	clientHelloMessage := &ClientHelloMessage{}
	err := msgpack.Unmarshal(rawBytes, clientHelloMessage)
	if err != nil {
		return nil, err
	}
	return clientHelloMessage, nil
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
	initiatorConnected *bool,
	responderAddresses *[]Address,
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

func NewDisconnectedMessage(id Address) DisconnectedMessage {
	return DisconnectedMessage{
		Message: Message{
			Type: Disconnected,
		},
		ID: id,
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
