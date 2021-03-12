package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/domain/dtos"
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
	initiatorsPublicKeyHex := path.Base(r.URL.Path)
	initiatorsPublicKey, err := dtos.FromHex(initiatorsPublicKeyHex)
	if err != nil {
		log.Error("fromHex:", err)
		return
	}
	connection, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	defer connection.Close()
	log.Infof("Remote address: %s", connection.RemoteAddr().String())

	client, err := c.saltyRTCServer.OnClientConnect(
		*initiatorsPublicKey,
		connection,
	)
	if err != nil {
		log.Error("onClientConnect:", err)
		return
	}

	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}
		log.Info("received new message")
		err = c.saltyRTCServer.OnMessage(*initiatorsPublicKey, client, message)
		if err != nil {
			log.Error("onClientConnect:", err)
			return
		}
	}
}
