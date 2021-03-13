package values

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
	YourCookie   Key      `msgpack:"your_cookie"`
	SubProtocols []string `msgpack:"subprotocols"`
	PingInterval int      `msgpack:"ping_interval"`
	YourKey      Key      `msgpack:"your_key"`
}

type ServerAuthMessage struct {
	Message
	YourCookie         Key       `msgpack:"your_cookie"`
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
