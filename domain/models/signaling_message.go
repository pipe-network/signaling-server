package models

import (
	"encoding/json"
	"fmt"
	"github.com/pipe-network/signaling-server/domain/values"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/nacl/box"
)

const SignalingMessageMinByteLength = 25

type SignalingMessage struct {
	Nonce values.Nonce
	Data  values.TypedMessage
}

func NewSignalingMessage(nonce values.Nonce, data values.TypedMessage) SignalingMessage {
	return SignalingMessage{
		Nonce: nonce,
		Data:  data,
	}
}

func (m SignalingMessage) DataBytes() ([]byte, error) {
	dataBytes, err := msgpack.Marshal(m.Data)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}

func (m SignalingMessage) Bytes() ([]byte, error) {
	var bytes []byte
	nonceBytes := m.Nonce.Bytes()
	bytes = append(bytes, nonceBytes[:]...)

	dataBytes, err := m.DataBytes()
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, dataBytes...)
	return bytes, nil
}

func (m SignalingMessage) EncryptBytes(publicKey, privateKey values.Key) ([]byte, error) {
	nonceBytes := m.Nonce.Bytes()
	dataBytes, err := m.DataBytes()
	if err != nil {
		return nil, err
	}
	publicKeyBytes := publicKey.Bytes()
	privateKeyBytes := privateKey.Bytes()

	encryptedDataBytes := box.Seal(nonceBytes[:], dataBytes, &nonceBytes, &publicKeyBytes, &privateKeyBytes)
	return encryptedDataBytes, nil
}

func (m SignalingMessage) String() string {
	data, _ := json.Marshal(m.Data)
	return fmt.Sprintf("Nonce:\n%s\nData:\n%s", m.Nonce.String(), string(data))
}
