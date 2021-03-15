package application

import (
	"github.com/pipe-network/signaling-server/interface/controllers"
	"log"
	"net/http"
)

type ServerAddress string

type MainApplication struct {
	SignallingController controllers.SignalingController

	serverAddress ServerAddress
}

func NewMainApplication(
	signallingController controllers.SignalingController,
	serverAddress ServerAddress,
) MainApplication {
	return MainApplication{
		SignallingController: signallingController,
		serverAddress:        serverAddress,
	}
}

func (a *MainApplication) Run() {
	log.SetFlags(0)
	http.HandleFunc("/", a.SignallingController.WebSocket)
	log.Printf("Running on: https://%s", a.serverAddress)
	log.Fatal(
		http.ListenAndServeTLS(
			string(a.serverAddress),
			"cert.pem",
			"key.pem",
			nil,
		),
	)
}
