package application

import (
	"flag"
	"github.com/pipe-network/signaling-server/src/interface/controllers"
	"log"
	"net/http"
)

type MainApplication struct {
	SignallingController controllers.SignalingController
}

func NewMainApplication(
	signallingController controllers.SignalingController,
) MainApplication {
	return MainApplication{
		SignallingController: signallingController,
	}
}

func (a *MainApplication) Run() {
	address := a.parseFlags()
	log.SetFlags(0)
	http.HandleFunc("/", a.SignallingController.WebSocket)
	log.Printf("Running on: http://%s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func (a MainApplication) parseFlags() string {
	address := flag.String("address", "localhost:8080", "http service address")
	if address == nil {
		return ""
	}
	flag.Parse()
	return *address
}
