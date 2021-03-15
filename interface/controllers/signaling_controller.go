package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/domain/values"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
)

type SignalingController struct {
	upgrader       websocket.Upgrader
	saltyRTCServer *services.SaltyRTCService
}

func NewSignalingController(
	upgrader websocket.Upgrader,
	saltyRTCService *services.SaltyRTCService,
) SignalingController {
	return SignalingController{
		upgrader:       upgrader,
		saltyRTCServer: saltyRTCService,
	}
}

func (c *SignalingController) WebSocket(w http.ResponseWriter, r *http.Request) {
	var err error
	log.Info("New request at ", r.RequestURI)
	initiatorsPublicKeyHex := path.Base(r.URL.Path)
	initiatorsPublicKey, err := values.FromHex(initiatorsPublicKeyHex)
	if err != nil {
		log.Errorf("fromHex: %v", err)
		return
	}
	connection, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade: %v", err)
		return
	}
	defer connection.Close()
	client, err := c.saltyRTCServer.OnClientConnect(
		*initiatorsPublicKey,
		connection,
	)
	if err != nil {
		log.Errorf("onClientConnect: %v", err)
		client.DropConnection(values.InternalErrorCode)
		return
	}

	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Errorf("readmessage: %v", err)
			client.DropConnection(values.InternalErrorCode)
			return
		}
		err = c.saltyRTCServer.OnMessage(*initiatorsPublicKey, client, message)
		if err != nil {
			log.Errorf("onmessage: %v", err)
			client.DropConnection(values.InternalErrorCode)
			return
		}
	}
}
