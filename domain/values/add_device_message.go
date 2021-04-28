package values

import (
	"errors"
)

var (
	CannotSplitMessage = errors.New("cannot split the message, probably no nonce sent")
)

type AddDeviceMessage struct {
	PublicKey Key
	Nonce     [NonceByteLength]byte
	Data      []byte
}

func AddDeviceMessageFromBytes(message []byte) (*AddDeviceMessage, error) {
	var strippedMessage []byte

	if len(message) < NonceByteLength+KeyByteLength {
		return nil, CannotSplitMessage
	}

	publicKey := [KeyByteLength]byte{}
	nonce := [NonceByteLength]byte{}

	arrayPointer := 0
	copy(publicKey[:], message[arrayPointer:arrayPointer+KeyByteLength])
	arrayPointer += KeyByteLength
	copy(nonce[:], message[arrayPointer:arrayPointer+NonceByteLength])
	arrayPointer += NonceByteLength
	strippedMessage = append(strippedMessage, message[arrayPointer:]...)
	return &AddDeviceMessage{
		PublicKey: publicKey,
		Nonce:     nonce,
		Data:      strippedMessage,
	}, nil
}

func (m AddDeviceMessage) ToBytes() []byte {
	var packedMessage []byte

	packedMessage = append(packedMessage, m.PublicKey[:]...)
	packedMessage = append(packedMessage, m.Nonce[:]...)
	packedMessage = append(packedMessage, m.Data...)
	return packedMessage
}
