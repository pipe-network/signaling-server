package dtos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/nacl/box"
)

type SignalingMessage struct {
	Nonce Nonce
	Data  interface{}
}

func NewSignalingMessage(destinationAddress Address, data interface{}) (*SignalingMessage, error) {
	combinedSequenceNumber, err := NewSequenceNumber()
	if err != nil {
		return nil, err
	}

	cookie, err := NewRandomCookie()
	if err != nil {
		return nil, err
	}

	return &SignalingMessage{
		Nonce: Nonce{
			Cookie:         *cookie,
			Source:         Server,
			Destination:    destinationAddress,
			SequenceNumber: *combinedSequenceNumber,
			OverflowNumber: NewOverflowNumber(),
		},
		Data: data,
	}, nil
}

func SignalingMessageFromBytes(bytes []byte, publicKey, privateKey Key) (*SignalingMessage, error) {
	var nonceBytes [24]byte
	var encodedDataBytes []byte
	var unknownMessage interface{}

	publicKeyBytes := publicKey.Bytes()
	privateKeyBytes := privateKey.Bytes()

	copy(nonceBytes[:], bytes[:NonceByteLength])
	encodedDataBytes = bytes[NonceByteLength:]

	decodedDataBytes, ok := box.Open(nil, encodedDataBytes, &nonceBytes, &publicKeyBytes, &privateKeyBytes)
	if !ok {
		return nil, errors.New("decryption failed")
	}

	err := msgpack.Unmarshal(decodedDataBytes, &unknownMessage)
	if err != nil {
		return nil, err
	}

	// messageType := unknownMessage.(map[string]interface{})["type"]

	return &SignalingMessage{
		Nonce: NonceFromBytes(nonceBytes),
		Data:  unknownMessage,
	}, nil
}

func (m SignalingMessage) Bytes() ([]byte, error) {
	var bytes []byte
	bytes = append(bytes, m.Nonce.Bytes()...)

	dataBytes, err := msgpack.Marshal(m.Data)
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, dataBytes...)
	return bytes, nil
}

func (m SignalingMessage) String() string {
	data, _ := json.Marshal(m.Data)
	return fmt.Sprintf("Nonce:\n%s\nData:\n%s", m.Nonce.String(), string(data))
}
