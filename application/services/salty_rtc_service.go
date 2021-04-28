package services

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/domain/models"
	"github.com/pipe-network/signaling-server/domain/values"
	log "github.com/sirupsen/logrus"
	"time"
)

const MinPingInterval = 0

var (
	InvalidSubProtocols = errors.New("invalid subprotocols")
	InvalidPingInterval = errors.New("invalid ping interval, shall be greater than 0")
	InvalidKey          = errors.New("invalid key")
	NoRoomInitiated     = errors.New("no room was initiated")
)

type SaltyRTCService interface {
	OnClientConnect(initiatorsPublicKey values.Key, connection *websocket.Conn) (*models.Client, error)
	OnMessage(initiatorsPublicKey values.Key, client *models.Client, message []byte) error
}

type SaltyRTCServiceImpl struct {
	rooms *models.Rooms

	keyPairStorage        ports.KeyPairStorage
	notificationService   ports.NotificationService
	deviceTokenRepository ports.DeviceTokenRepository
}

func NewSaltyRTCServiceImpl(
	keyPairStorage ports.KeyPairStorage,
	notificationService ports.NotificationService,
	deviceTokenRepository ports.DeviceTokenRepository,
) *SaltyRTCServiceImpl {
	return &SaltyRTCServiceImpl{
		rooms:                 models.NewRooms(),
		keyPairStorage:        keyPairStorage,
		notificationService:   notificationService,
		deviceTokenRepository: deviceTokenRepository,
	}
}

func (s *SaltyRTCServiceImpl) OnClientConnect(
	initiatorsPublicKey values.Key,
	connection *websocket.Conn,
) (*models.Client, error) {
	room := s.rooms.GetOrCreateRoom(initiatorsPublicKey)

	log.Infof("Rooms: %d", s.rooms.Size())

	client, err := models.NewClient(connection, room)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}

	connection.SetCloseHandler(func(code int, text string) error {
		s.cleanup(client, room)
		return nil
	})

	room.AddClient(client)
	serverHelloMessage := values.NewServerHelloMessage(client.SessionPublicKey)
	signalingMessage := models.NewSignalingMessage(client.Nonce(), &serverHelloMessage)
	signalingMessageBytes, err := signalingMessage.Bytes()
	if err != nil {
		log.Errorf("dropping connection: could not get bytes of signaling message: %v", err)
		client.DropConnection(values.InternalErrorCode)
		s.cleanup(client, room)
		return nil, err
	}
	err = client.SendBytes(signalingMessageBytes)
	if err != nil {
		log.Errorf("dropping connection: could not send bytes with signaling message: %v", err)
		client.DropConnection(values.InternalErrorCode)
		s.cleanup(client, room)
		return nil, err
	}
	return client, nil
}

func (s *SaltyRTCServiceImpl) ReadMessageLoop(initiatorsPublicKey values.Key, client *models.Client) {
	room := s.rooms.GetOrCreateRoom(initiatorsPublicKey)

	for {
		_, message, err := client.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, values.HandoverOfTheSignalingChannelCode.Int()) {
				log.Infof("Handover of %d", client.Address)
				err = client.CloseConnection()
				if err != nil {
					log.Errorf("closing connection: could not close connection: %v", err)
				}
				s.cleanup(client, room)
				return
			}
			log.Errorf("dropping connection: could read client message: %v", err)
			client.DropConnection(values.InternalErrorCode)
			s.cleanup(client, room)
			return
		}
		err = s.OnMessage(initiatorsPublicKey, client, message)
		if err != nil {
			log.Errorf("dropping connection: could not process onmessage: %v", err)
			client.DropConnection(values.InternalErrorCode)
			s.cleanup(client, room)
			return
		}
	}
}

func (s *SaltyRTCServiceImpl) OnMessage(initiatorsPublicKey values.Key, client *models.Client, message []byte) error {
	nonce, dataBytes, err := s.splitMessage(message)
	if err != nil {
		return err
	}

	room := s.rooms.GetRoom(initiatorsPublicKey)
	if room == nil {
		return NoRoomInitiated
	}

	if nonce.Destination == values.ServerAddress {
		if client.IncomingNonceEmpty() {
			client.SetIncomingNonce(nonce)
		}

		err = client.ValidateNonce(nonce)
		if err != nil {
			return err
		}

		// Increment incoming combined sequence number for next request after validation
		err = client.IncrementIncomingCombinedSequenceNumber()
		if err != nil {
			return err
		}

		// Client is is initiator and already authenticated, so it's probably the drop responder message
		if client.IsAuthenticated() && client.IsInitiator() {
			dropResponderMessage, err := values.DecodeDropResponderMessageFromBytes(
				dataBytes,
				nonce.Bytes(),
				initiatorsPublicKey,
				client.SessionPrivateKey,
			)
			if err != nil {
				return err
			}
			return s.onDropResponderMessage(*dropResponderMessage, room)
		}

		clientsPermanentPublicKey := values.Key{}
		if !client.PermanentPublicKey.Empty() {
			clientsPermanentPublicKey = client.PermanentPublicKey
		} else {
			clientsPermanentPublicKey = initiatorsPublicKey
		}
		clientAuthMessage, err := values.DecodeClientAuthMessageFromBytes(
			dataBytes,
			nonce.Bytes(),
			clientsPermanentPublicKey,
			client.SessionPrivateKey,
		)

		// If the decryption failed, check if it's a client-hello message, otherwise throw err
		if err == values.DecryptionFailed {
			clientHelloMessage, err := values.DecodeClientHelloMessageFromBytes(dataBytes)
			if err != nil {
				return errors.New("cannot unpack neither client-hello nor client-auth")
			}

			err = s.onClientHelloMessage(client, *clientHelloMessage)
			if err != nil {
				return err
			}

			return nil
		} else if err != nil {
			client.DropConnection(values.ProtocolErrorCode)
			s.cleanup(client, room)
			return err
		}

		err = s.onClientAuthMessage(client, room, *clientAuthMessage)
		if err != nil {
			return err
		}

		return nil
	} else {
		toClient := room.Client(nonce.Destination)
		if (client.IsInitiator() && toClient.IsResponder() || client.IsResponder() && toClient.IsInitiator()) && client.IsAuthenticated() && toClient.IsAuthenticated() {
			err := toClient.SendBytes(message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SaltyRTCServiceImpl) cleanup(client *models.Client, room *models.Room) {
	log.Debug("Cleanup after close of: ", client.Address)
	s.broadcastDisconnected(room, client)
	if client.IsResponder() {
		room.ReleaseAddress(client.Address)
	}
	room.RemoveClient(client)
	client.Flush()
}

func (s *SaltyRTCServiceImpl) splitMessage(message []byte) (values.Nonce, []byte, error) {
	messageLength := len(message)
	if messageLength < models.SignalingMessageMinByteLength {
		return values.Nonce{}, nil, errors.New(fmt.Sprintf("message too short %v bytes", messageLength))
	}

	var nonceBytes [values.NonceByteLength]byte
	var dataBytes []byte
	copy(nonceBytes[:], message[:values.NonceByteLength])
	dataBytes = message[values.NonceByteLength:]

	nonce := values.NonceFromBytes(nonceBytes)
	return nonce, dataBytes, nil
}

func (s *SaltyRTCServiceImpl) onClientAuthMessage(
	client *models.Client,
	room *models.Room,
	clientAuthMessage values.ClientAuthMessage,
) error {
	if !client.OutgoingCookie.Equal(clientAuthMessage.YourCookie) {
		return models.InvalidCookie
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
			s.cleanup(client, room)
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
			s.cleanup(client, room)
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

	outgoingNonce := client.Nonce()
	var initiatorConnected *bool
	var responderAddresses *[]values.Address

	if client.IsInitiator() {
		roomResponders := room.Responders()
		responderAddressesTemp := make([]values.Address, len(roomResponders))
		for i, responder := range roomResponders {
			responderAddressesTemp[i] = responder.Address
		}

		responderAddresses = &responderAddressesTemp
	} else {
		initiatorConnectedTemp := room.Initiator() != nil
		initiatorConnected = &initiatorConnectedTemp
		if room.Initiator() == nil {
			device, err := s.deviceTokenRepository.DeviceByPublicKey(room.InitiatorsPublicKey.HexString())
			if err != nil {
				return err
			}

			err = s.notificationService.Notify(
				"Pipe Network is syncing",
				"Pipe Network is syncing for you",
				map[string]string{
					"type":      "wakeup",
					"publicKey": client.PermanentPublicKey.HexString(),
				},
				device.Token,
			)
			if err != nil {
				panic(err)
			}
		}
	}

	serverAuthMessage := values.NewServerAuthMessage(
		client.IncomingCookie,
		client.SessionPublicKey,
		client.PermanentPublicKey,
		s.keyPairStorage.PrivateKey(),
		outgoingNonce,
		initiatorConnected,
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

func (s *SaltyRTCServiceImpl) onClientHelloMessage(
	client *models.Client,
	clientHelloMessage values.ClientHelloMessage,
) error {
	log.Infof("Setting client permanent public key: %s", clientHelloMessage.Key.HexString())
	client.SetPermanentPublicKey(clientHelloMessage.Key)
	return nil
}

func (s *SaltyRTCServiceImpl) broadcastNewInitiatorMessage(room *models.Room) error {
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

func (s *SaltyRTCServiceImpl) onDropResponderMessage(
	dropResponderMessage values.DropResponderMessage,
	room *models.Room,
) error {
	client := room.Client(dropResponderMessage.ID)
	if client != nil {
		reason := dropResponderMessage.Reason
		if reason == values.CloseCode(0) {
			reason = values.DroppedByInitiatorCode
		}

		client.DropConnection(reason)
		s.cleanup(client, room)
	}
	return nil
}

func (s *SaltyRTCServiceImpl) broadcastNewResponderMessage(responderClient *models.Client, room *models.Room) error {
	newResponderMessage := values.NewNewResponderMessage(responderClient.Address)
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

func (s *SaltyRTCServiceImpl) broadcastDisconnected(room *models.Room, disconnectedClient *models.Client) {
	if disconnectedClient.IsAuthenticated() {
		disconnectedMessage := values.NewDisconnectedMessage(disconnectedClient.Address)
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
