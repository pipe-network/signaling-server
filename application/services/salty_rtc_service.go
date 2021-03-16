package services

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/domain/models"
	"github.com/pipe-network/signaling-server/domain/values"
	"time"
)

const MinPingInterval = 0

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
	InvalidSubProtocols   = errors.New("invalid subprotocols")
	InvalidPingInterval   = errors.New("invalid ping interval, shall be greater than 0")
	InvalidKey            = errors.New("invalid key")
	NoRoomInitiated       = errors.New("no room was initiated")
)

type ISaltyRTCService interface {
	OnClientConnect(initiatorsPublicKey values.Key, connection *websocket.Conn) (*models.Client, error)
	OnMessage(initiatorsPublicKey values.Key, client *models.Client, message []byte) error
}

type SaltyRTCService struct {
	rooms *models.Rooms

	signalingMessageService ISignalingMessageService
	keyPairStorage          ports.KeyPairStoragePort
}

func NewSaltyRTCService(
	signalingMessageService ISignalingMessageService,
	keyPairStorage ports.KeyPairStoragePort,
) *SaltyRTCService {
	return &SaltyRTCService{
		rooms:                   models.NewRooms(),
		signalingMessageService: signalingMessageService,
		keyPairStorage:          keyPairStorage,
	}
}

func (s *SaltyRTCService) OnClientConnect(
	initiatorsPublicKey values.Key,
	connection *websocket.Conn,
) (*models.Client, error) {
	var room *models.Room
	room = s.rooms.GetRoom(initiatorsPublicKey)
	if room == nil {
		room = models.NewRoom(initiatorsPublicKey)
		s.rooms.AddRoom(room)
	}

	client, err := models.NewClient(connection)
	connection.SetCloseHandler(func(code int, text string) error {
		s.broadcastDisconnected(room, client)
		if client.IsResponder() {
			room.ReleaseAddress(client.Address)
		}
		room.RemoveClient(client)
		return nil
	})
	if err != nil {
		return nil, err
	}
	room.AddClient(client)

	serverHelloMessage := values.NewServerHelloMessage(client.SessionPublicKey)
	signalingMessage := models.NewSignalingMessage(client.Nonce(), &serverHelloMessage)
	signalingMessageBytes, err := signalingMessage.Bytes()
	if err != nil {
		return nil, err
	}
	err = client.SendBytes(signalingMessageBytes)
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
	if client.IncomingNonceEmpty() {
		client.SetIncomingNonce(nonce)
	}
	err := s.validateNonce(nonce, client)
	if err != nil {
		return err
	}

	// Increment incoming combined sequence number for next request
	err = client.IncrementIncomingCombinedSequenceNumber()
	if err != nil {
		return err
	}

	room := s.rooms.GetRoom(initiatorsPublicKey)
	if room == nil {
		return NoRoomInitiated
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
			err = s.onClientHelloMessage(client, *clientHelloMessage)
			if err != nil {
				return err
			}
		} else if err != nil {
			client.DropConnection(values.ProtocolErrorCode)
			return err
		} else {
			err = s.onClientAuthMessage(client, room, *clientAuthMessage)
			if err != nil {
				return err
			}
		}

	} else {
		toClient := room.Client(nonce.Destination)
		err := toClient.SendBytes(message)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *SaltyRTCService) validateNonce(nonce values.Nonce, client *models.Client) error {
	isAddressedToServer := nonce.Destination == values.ServerAddress
	combinedSequenceNumber := values.NewCombinedSequenceNumber(nonce.SequenceNumber, nonce.OverflowNumber)

	// Validate destination address
	if !isAddressedToServer && !client.IsP2PAllowed(nonce.Destination) {
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

		if !client.IsCombinedSequenceNumberValid(combinedSequenceNumber) {
			return InvalidSequenceNumber
		}
	}
	return nil
}

func (s *SaltyRTCService) onClientAuthMessage(
	client *models.Client,
	room *models.Room,
	clientAuthMessage values.ClientAuthMessage,
) error {
	if !client.OutgoingCookie.Equal(clientAuthMessage.YourCookie) {
		return InvalidCookie
	}

	if !clientAuthMessage.ContainsSubProtocol(values.SaltyRTCSubprotocol) {
		return InvalidSubProtocols
	}

	if clientAuthMessage.PingInterval < MinPingInterval {
		return InvalidPingInterval
	} else if clientAuthMessage.PingInterval > MinPingInterval {
		pingPeriod, _ := time.ParseDuration(fmt.Sprintf("%ds", clientAuthMessage.PingInterval))
		go client.PingTicker(pingPeriod, models.DefaultPongWait)
	}

	if !clientAuthMessage.YourKey.Empty() {
		if !clientAuthMessage.YourKey.Equals(s.keyPairStorage.PublicKey()) {
			client.DropConnection(values.InvalidKeyCode)
			return InvalidKey
		}
	}

	if client.PermanentPublicKey.Empty() {
		client.SetPermanentPublicKey(room.InitiatorsPublicKey)
		room.KickCurrentInitiator()
		client.AssignToInitiator()
		err := s.broadcastNewInitiatorMessage(room)
		if err != nil {
			return err
		}
	} else {
		nextFreeResponderAddress, err := room.NextFreeResponderAddress()
		if err != nil {
			client.DropConnection(values.PathFullCode)
			return err
		}
		client.SetAddress(*nextFreeResponderAddress)
		room.ReserveAddress(*nextFreeResponderAddress)
		err = s.broadcastNewResponderMessage(client, room)
		if err != nil {
			return err
		}
	}

	client.MarkAsAuthenticated()

	responderAddresses := make([]int, 0)
	for _, responder := range room.Responders() {
		responderAddresses = append(responderAddresses, int(responder.Address))
	}

	outgoingNonce := client.Nonce()
	serverAuthMessage := values.NewServerAuthMessage(
		client.IncomingCookie,
		client.SessionPublicKey,
		client.PermanentPublicKey,
		s.keyPairStorage.PrivateKey(),
		outgoingNonce,
		room.Initiator() != nil,
		responderAddresses,
	)
	signalingMessage := models.NewSignalingMessage(outgoingNonce, &serverAuthMessage)
	encryptedSignalingMessageBytes, err := signalingMessage.EncryptBytes(
		client.PermanentPublicKey,
		client.SessionPrivateKey,
	)
	if err != nil {
		return err
	}
	err = client.SendBytes(encryptedSignalingMessageBytes)
	if err != nil {

		return err
	}
	return nil
}

func (s *SaltyRTCService) onClientHelloMessage(
	client *models.Client,
	clientHelloMessage values.ClientHelloMessage,
) error {
	client.SetPermanentPublicKey(clientHelloMessage.Key)
	return nil
}

func (s *SaltyRTCService) broadcastNewInitiatorMessage(room *models.Room) error {
	newInitiatorMessage := values.NewNewInitiatorMessage()
	responders := room.Responders()

	for _, responderClient := range responders {
		signalingMessage := models.NewSignalingMessage(responderClient.Nonce(), &newInitiatorMessage)
		signalingMessageBytes, err := signalingMessage.EncryptBytes(
			responderClient.PermanentPublicKey,
			responderClient.SessionPrivateKey,
		)
		if err != nil {
			return err
		}
		_ = responderClient.SendBytes(signalingMessageBytes)
	}
	return nil
}

func (s *SaltyRTCService) broadcastNewResponderMessage(responderClient *models.Client, room *models.Room) error {
	newResponderMessage := values.NewNewResponderMessage(int(responderClient.Address))
	initiator := room.Initiator()
	if initiator == nil {
		return nil
	}
	signalingMessage := models.NewSignalingMessage(initiator.Nonce(), &newResponderMessage)
	signalingMessageBytes, err := signalingMessage.EncryptBytes(
		initiator.PermanentPublicKey,
		initiator.SessionPrivateKey,
	)
	if err != nil {
		return err
	}
	err = initiator.SendBytes(signalingMessageBytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *SaltyRTCService) broadcastDisconnected(room *models.Room, disconnectedClient *models.Client) {
	if disconnectedClient.IsAuthenticated() {
		disconnectedMessage := values.NewDisconnectedMessage(int(disconnectedClient.Address))
		if disconnectedClient.IsInitiator() {
			for _, responderClient := range room.Responders() {
				signalingMessage := models.NewSignalingMessage(responderClient.Nonce(), &disconnectedMessage)
				signalingMessageBytes, _ := signalingMessage.EncryptBytes(
					responderClient.PermanentPublicKey,
					responderClient.SessionPrivateKey,
				)
				_ = responderClient.SendBytes(signalingMessageBytes)
			}
		}
		if disconnectedClient.IsResponder() {
			initiatorClient := room.Initiator()
			if initiatorClient == nil {
				return
			}
			signalingMessage := models.NewSignalingMessage(initiatorClient.Nonce(), &disconnectedMessage)
			signalingMessageBytes, _ := signalingMessage.EncryptBytes(
				initiatorClient.PermanentPublicKey,
				initiatorClient.SessionPrivateKey,
			)
			_ = initiatorClient.SendBytes(signalingMessageBytes)
		}
	}
}
