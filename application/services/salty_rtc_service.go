package services

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/models"
	"github.com/pipe-network/signaling-server/domain/values"
	log "github.com/sirupsen/logrus"
)

var (
	NotAllowedToRelay = func(destinationAddress values.Address) error {
		return errors.New(fmt.Sprintf("not allowed to relay messages to %x", destinationAddress))
	}
	IdentitiesDoNotMatch = func(clientAddress, sourceAddress values.Address) error {
		return errors.New(
			fmt.Sprintf("identities do not match, expected %x, got %x", clientAddress, sourceAddress),
		)
	}
	InvalidCookie         = errors.New("invalid cookie")
	InvalidSequenceNumber = errors.New("invalid sequence number")
)

type ISaltyRTCService interface {
	OnClientConnect(initiatorsPublicKey values.Key, connection *websocket.Conn) (*models.Client, error)
	OnMessage(initiatorsPublicKey values.Key, client *models.Client, message []byte) error
}

type SaltyRTCService struct {
	rooms *models.Rooms

	signalingMessageService ISignalingMessageService
}

func NewSaltyRTCService(signalingMessageService ISignalingMessageService) *SaltyRTCService {
	return &SaltyRTCService{
		rooms:                   models.NewRooms(),
		signalingMessageService: signalingMessageService,
	}
}

func (s *SaltyRTCService) OnClientConnect(
	initiatorsPublicKey values.Key,
	connection *websocket.Conn,
) (*models.Client, error) {
	room := models.NewRoom(initiatorsPublicKey)
	client, err := models.NewClient(connection)
	if err != nil {
		return nil, err
	}

	s.rooms.AddRoom(room)
	room.AddClient(client)

	serverHelloMessage := values.NewServerHelloMessage(client.SessionPublicKey)
	signalingMessage := models.NewSignalingMessage(client, &serverHelloMessage)
	log.Info("OutgoingSignalingMessage:", signalingMessage.String())
	err = client.SendSignalingMessage(signalingMessage)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *SaltyRTCService) OnMessage(initiatorsPublicKey values.Key, client *models.Client, message []byte) error {
	messageLength := len(message)
	if messageLength < models.SignalingMessageMinByteLength {
		return errors.New(fmt.Sprintf("message too short %v bytes", messageLength))
	}

	var nonceBytes [values.NonceByteLength]byte
	var dataBytes []byte
	copy(nonceBytes[:], message[:values.NonceByteLength])
	dataBytes = message[values.NonceByteLength:]

	nonce := s.signalingMessageService.NonceFromBytes(nonceBytes)
	err := s.validateNonce(nonce, client)
	if err != nil {
		return err
	}

	if nonce.Destination == values.ServerAddress {
		clientAuthMessage, err := s.signalingMessageService.DecodeClientAuthMessageFromBytes(
			dataBytes,
			nonceBytes,
			initiatorsPublicKey,
			client.SessionPrivateKey,
		)

		// If the decryption failed, check if it's a client-hello message, otherwise throw err
		if err == DecryptionFailed {
			clientHelloMessage, err := s.signalingMessageService.ClientHelloMessageFromBytes(dataBytes)
			if err != nil {
				return errors.New("cannot unpack neither client-hello nor client-auth")
			}

			s.onClientHelloMessage(client, nonce, *clientHelloMessage)
		} else if err != nil {
			return err
		} else {
			s.onClientAuthMessage(client, nonce, *clientAuthMessage)
		}

	} else {
		// Relay message to destination address
		// source peer tries to communicate with destination peer
	}

	return nil

}

func (s *SaltyRTCService) validateNonce(nonce values.Nonce, client *models.Client) error {
	isAddressedToServer := nonce.Destination == values.ServerAddress

	// Validate destination address
	if !isAddressedToServer && client.IsP2PAllowed(nonce.Destination) {
		return NotAllowedToRelay(nonce.Destination)
	}

	// Validate source address
	if nonce.Source != client.Address {
		return IdentitiesDoNotMatch(client.Address, nonce.Source)
	}

	if isAddressedToServer {
		if !client.IsCookieValid(nonce.Cookie) {
			return InvalidCookie
		}
		if !client.IsCombinedSequenceNumberValid(nonce.SequenceNumber, nonce.OverflowNumber) {
			return InvalidSequenceNumber
		}
		err := client.IncrementIncomingCombinedSequenceNumber()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SaltyRTCService) onClientAuthMessage(
	client *models.Client,
	nonce values.Nonce,
	clientAuthMessage values.ClientAuthMessage,
) {
	log.Info("it's an onClientAuthMessage\n")
	log.Info(nonce.String())
	log.Info(clientAuthMessage)
	// do it
}

func (s *SaltyRTCService) onClientHelloMessage(
	client *models.Client,
	nonce values.Nonce,
	clientHelloMessage values.ClientHelloMessage,
) {
	log.Info("it's an onClientHello\n")
	log.Info(nonce.String())
	log.Info(clientHelloMessage)
}
