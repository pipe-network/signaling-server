package messages

import "github.com/pipe-network/signaling-server/src/domain/dtos"

type ServerHelloMessage struct {
	MessageType dtos.MessageType `msgpack:"type"`
	Key         dtos.Key
}

func NewServerHelloMessage(sessionKey dtos.Key) ServerHelloMessage {
	return ServerHelloMessage{
		MessageType: dtos.ServerHello,
		Key:         sessionKey,
	}
}
