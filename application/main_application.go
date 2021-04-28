package application

import (
	"fmt"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/interface/controllers"
	"log"
	"net/http"
)

type MainApplication struct {
	signallingController controllers.SignalingController
	addDeviceController  controllers.AddDeviceController
	flagService          services.FlagService
}

func NewMainApplication(
	flagService services.FlagService,
	signallingController controllers.SignalingController,
	addDeviceController controllers.AddDeviceController,
) MainApplication {
	return MainApplication{
		signallingController: signallingController,
		addDeviceController:  addDeviceController,
		flagService:          flagService,
	}
}

func (a *MainApplication) Run() {
	address := a.flagService.String(services.Address)
	port := a.flagService.Int(services.Port)
	http.HandleFunc("/add-device-token", a.addDeviceController.Websocket)
	http.HandleFunc("/", a.signallingController.WebSocket)
	log.Printf("Running on: https://%s:%d", address, port)
	log.Fatal(
		http.ListenAndServeTLS(
			fmt.Sprintf("%s:%d", address, port),
			a.flagService.String(services.TLSCertFile),
			a.flagService.String(services.TLSKeyFile),
			nil,
		),
	)
}
