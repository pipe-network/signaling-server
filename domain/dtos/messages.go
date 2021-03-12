package dtos

type ServerHelloMessage struct {
	MessageType MessageType `msgpack:"type"`
	Key         Key         `msgpack:"key"`
}

func NewServerHelloMessage(sessionPublicKey Key) ServerHelloMessage {
	return ServerHelloMessage{
		MessageType: ServerHello,
		Key:         sessionPublicKey,
	}
}
