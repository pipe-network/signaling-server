package models

import (
	"encoding/json"
	"fmt"
	"github.com/pipe-network/signaling-server/domain/values"
	"github.com/vmihailenco/msgpack/v5"
)

const SignalingMessageMinByteLength = 25

type SignalingMessage struct {
	Nonce values.Nonce
	Data  values.TypedMessage
}

func NewSignalingMessage(client *Client, data values.TypedMessage) SignalingMessage {
	return SignalingMessage{
		Nonce: values.Nonce{
			Cookie:         client.Cookie,
			Source:         values.ServerAddress,
			Destination:    client.Address,
			SequenceNumber: client.SequenceNumber,
			OverflowNumber: client.OverflowNumber,
		},
		Data: data,
	}
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
