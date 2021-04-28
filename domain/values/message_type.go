package values

const (
	ServerHello      MessageType = "server-hello"
	ClientHello      MessageType = "client-hello"
	ClientAuth       MessageType = "client-auth"
	ServerAuth       MessageType = "server-auth"
	NewInitiator     MessageType = "new-initiator"
	NewResponder     MessageType = "new-responder"
	DropResponder    MessageType = "drop-responder"
	Disconnected     MessageType = "disconnected"
	SendError        MessageType = "send-error"
	AddDeviceControl MessageType = "add-device-control"
)

type MessageType string
