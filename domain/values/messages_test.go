package values

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/nacl/box"
	"testing"
)

func TestClientAuthMessage_ContainsSubProtocol(t *testing.T) {
	clientAuthMessage := ClientAuthMessage{
		SubProtocols: []string{
			"123",
			"456",
		},
	}
	assert.True(t, clientAuthMessage.ContainsSubProtocol("123"))
	assert.True(t, clientAuthMessage.ContainsSubProtocol("456"))
	assert.False(t, clientAuthMessage.ContainsSubProtocol("678"))
}

func TestNewServerAuthMessage(t *testing.T) {
	cookie := Cookie{0x1, 0x2}
	responderAddresses := []Address{2}
	outgoingSessionPublicKey := Key{0x1}
	clientPermanentPublicKey := Key{0x2}
	serverPermanentPrivateKey := Key{0x3}
	peersPublicKeyBytes := clientPermanentPublicKey.Bytes()
	privateKeyBytes := serverPermanentPrivateKey.Bytes()
	nonce := Nonce{
		Cookie:         Cookie{0x3},
		Source:         0,
		Destination:    10,
		OverflowNumber: OverflowNumber{0},
		SequenceNumber: SequenceNumber{0},
	}

	nonceBytes := nonce.Bytes()
	var signedKeys []byte
	signedKeys = append(signedKeys, outgoingSessionPublicKey[:]...)
	signedKeys = append(signedKeys, clientPermanentPublicKey[:]...)
	signedKeys = box.Seal(nil, signedKeys, &nonceBytes, &peersPublicKeyBytes, &privateKeyBytes)
	initiatorConnect := false

	serverAuthMessage := NewServerAuthMessage(
		cookie,
		outgoingSessionPublicKey,
		clientPermanentPublicKey,
		serverPermanentPrivateKey,
		nonce,
		&initiatorConnect,
		&responderAddresses,
	)
	assert.Equal(
		t,
		ServerAuthMessage{
			Message: Message{
				Type: ServerAuth,
			},
			YourCookie:         cookie,
			SignedKeys:         signedKeys,
			InitiatorConnected: &initiatorConnect,
			Responders:         &responderAddresses,
		},
		serverAuthMessage,
	)
}
