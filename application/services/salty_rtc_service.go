package services

import (
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/dtos"
	"github.com/pipe-network/signaling-server/domain/models"
	log "github.com/sirupsen/logrus"
)

type SaltyRTCService struct {
	rooms *models.Rooms
}

func NewSaltyRTCService() *SaltyRTCService {
	return &SaltyRTCService{
		rooms: models.NewRooms(),
	}
}

func (s *SaltyRTCService) OnClientConnect(
	initiatorsPublicKey dtos.Key,
	connection *websocket.Conn,
) (*models.Client, error) {
	room := models.NewRoom(initiatorsPublicKey)
	client, err := models.NewClient(connection)
	if err != nil {
		return nil, err
	}

	s.rooms.AddRoom(room)
	room.AddClient(client)

	serverHelloMessage := dtos.NewServerHelloMessage(client.SessionPublicKey)
	signalingMessage, err := dtos.NewSignalingMessage(client.Address, serverHelloMessage)
	if err != nil {
		return nil, err
	}
	log.Info("OutgoingSignalingMessage:", signalingMessage.String())
	err = client.SendSignalingMessage(*signalingMessage)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *SaltyRTCService) OnMessage(initiatorsPublicKey dtos.Key, client *models.Client, message []byte) error {
	signalingMessage, err := dtos.SignalingMessageFromBytes(message, initiatorsPublicKey, client.SessionPrivateKey)
	_ = s.rooms.GetRoom(initiatorsPublicKey)
	if err != nil {
		return err
	}

	log.Info("IncomingSignalingMessage:", signalingMessage.String())
	return nil

}
