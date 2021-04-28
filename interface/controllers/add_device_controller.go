package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/application/services"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AddDeviceController struct {
	upgrader         websocket.Upgrader
	addDeviceService services.AddDeviceService
}

func NewAddDeviceController(
	upgrader websocket.Upgrader,
	addDeviceService services.AddDeviceService,
) AddDeviceController {
	return AddDeviceController{
		upgrader:         upgrader,
		addDeviceService: addDeviceService,
	}
}

func (c *AddDeviceController) Websocket(w http.ResponseWriter, r *http.Request) {
	log.Info("New add device controller request")
	connection, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}
	defer func(connection *websocket.Conn) {
		err := connection.Close()
		if err != nil {
			log.Error(err)
		}
	}(connection)

	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Errorf("dropping connection: could read client message: %v", err)
			return
		}
		err = c.addDeviceService.OnAddDeviceMessage(connection, message)
		if err != nil {
			log.Error(err)
			err := connection.Close()
			if err != nil {
				log.Error(err)
			}
			return
		}
	}

}
