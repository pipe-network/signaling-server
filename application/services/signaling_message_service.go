package services

import (
	"errors"
	"github.com/pipe-network/signaling-server/domain/values"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/nacl/box"
)

var (
	DecryptionFailed error = errors.New("decryption failed")
)

type ISignalingMessageService interface {
	NonceFromBytes(bytes [values.NonceByteLength]byte) values.Nonce
	DecodeClientAuthMessageFromBytes(
		data []byte,
		nonce [values.NonceByteLength]byte,
		publicKey,
		privateKey values.Key,
	) (*values.ClientAuthMessage, error)
	ClientHelloMessageFromBytes(bytes []byte) (*values.ClientHelloMessage, error)
}

type SignalingMessageService struct{}

func NewSignalingMessageService() *SignalingMessageService {
	return &SignalingMessageService{}
}

func (s *SignalingMessageService) NonceFromBytes(bytes [values.NonceByteLength]byte) values.Nonce {
	var cookie values.Cookie
	var source, destination values.Address
	var overflowNumber values.OverflowNumber
	var sequenceNumber values.SequenceNumber
	var from int

	sourceBytes := make([]byte, 1)
	destinationBytes := make([]byte, 1)

	copy(cookie[:], bytes[from:len(cookie)])
	from += len(cookie)
	copy(sourceBytes[:], bytes[from:from+len(sourceBytes)])
	from += len(source.Bytes())
	copy(destinationBytes[:], bytes[from:from+len(destinationBytes)])
	from += len(destination.Bytes())
	copy(overflowNumber[:], bytes[from:from+len(overflowNumber)])
	from += len(overflowNumber)
	copy(sequenceNumber[:], bytes[from:from+len(sequenceNumber)])

	if len(sourceBytes) > 0 {
		source = values.Address(int(sourceBytes[0]))
	}

	if len(destinationBytes) > 0 {
		destination = values.Address(int(destinationBytes[0]))
	}

	return values.Nonce{
		Cookie:         cookie,
		Source:         source,
		Destination:    destination,
		OverflowNumber: overflowNumber,
		SequenceNumber: sequenceNumber,
	}
}

func (s *SignalingMessageService) DecodeClientAuthMessageFromBytes(
	data []byte,
	nonce [values.NonceByteLength]byte,
	publicKey,
	privateKey values.Key,
) (*values.ClientAuthMessage, error) {
	var clientAuthMessage values.ClientAuthMessage

	publicKeyBytes := publicKey.Bytes()
	privateKeyBytes := privateKey.Bytes()

	decodedDataBytes, ok := box.Open(nil, data, &nonce, &publicKeyBytes, &privateKeyBytes)
	if !ok {
		return nil, DecryptionFailed
	}

	err := msgpack.Unmarshal(decodedDataBytes, &clientAuthMessage)
	if err != nil {
		return nil, err
	}
	return &clientAuthMessage, nil
}

func (s *SignalingMessageService) ClientHelloMessageFromBytes(rawBytes []byte) (*values.ClientHelloMessage, error) {
	clientHelloMessage := &values.ClientHelloMessage{}
	err := msgpack.Unmarshal(rawBytes, clientHelloMessage)
	if err != nil {
		return nil, err
	}
	return clientHelloMessage, nil
}
