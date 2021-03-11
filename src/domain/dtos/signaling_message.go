package dtos

import "github.com/pipe-network/signaling-server/src/domain/dtos/messages"

type SignalingMessage struct {
	Nonce Nonce
	Data  interface{}
}

func NewServerHelloSignalingMessage(sessionKey Key, destination Address) (*SignalingMessage, error) {
	combinedSequenceNumber, err := NewCombinedSequenceNumber()
	if err != nil {
		return nil, err
	}

	cookie, err := NewRandomCookie()
	if err != nil {
		return nil, err
	}

	return &SignalingMessage{
		Nonce: Nonce{
			Cookie:                 *cookie,
			Source:                 Server,
			Destination:            destination,
			CombinedSequenceNumber: *combinedSequenceNumber,
		},
		Data: messages.NewServerHelloMessage(sessionKey),
	}, nil
}
