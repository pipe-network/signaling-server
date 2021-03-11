package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/src/domain/dtos"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/box"
	"net/http"
	"path"
)

type SignalingController struct {
	Upgrader websocket.Upgrader
}

func NewSignalingController(upgrader websocket.Upgrader) SignalingController {
	return SignalingController{
		Upgrader: upgrader,
	}
}

func (c *SignalingController) WebSocket(w http.ResponseWriter, r *http.Request) {
	initiatorsPublicKeyHex := path.Base(r.URL.Path)
	_, err := dtos.FromHex(initiatorsPublicKeyHex)
	if err != nil {
		log.Error(err)
		return
	}
	newPublicKey, newPrivateKey, err := box.GenerateKey(rand.Reader)
	// log.Info("Publickey: ", publicKeyHex)
	log.Info("NewPublicKey: ", hex.EncodeToString(newPublicKey[:]))
	log.Info("NewPrivateKey: ", hex.EncodeToString(newPrivateKey[:]))

	if err != nil {
		log.Error("generateKey:", err)
		return
	}

	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = connection.WriteMessage(messageType, message)
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
}
