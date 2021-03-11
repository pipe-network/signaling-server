package controllers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type SignalingController struct {
	Upgrader websocket.Upgrader
}

func NewSignalingController (upgrader websocket.Upgrader) SignalingController {
	return SignalingController{
		Upgrader: upgrader,
	}
}

func (c *SignalingController) WebSocket(w http.ResponseWriter, r *http.Request)  {
	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = connection.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
