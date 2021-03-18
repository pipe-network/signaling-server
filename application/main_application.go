package application

import (
	"github.com/pipe-network/signaling-server/interface/controllers"
	"log"
	"net/http"
)

type ServerAddress string
type TLSCertFilePath string
type TLSKeyFilePath string

type MainApplication struct {
	SignallingController controllers.SignalingController

	serverAddress   ServerAddress
	tlsCertFilePath TLSCertFilePath
	tlsKeyFilePath  TLSKeyFilePath
}

func NewMainApplication(
	signallingController controllers.SignalingController,
	serverAddress ServerAddress,
	tlsCertFilePath TLSCertFilePath,
	tlsKeyFilePath TLSKeyFilePath,
) MainApplication {
	return MainApplication{
		SignallingController: signallingController,
		serverAddress:        serverAddress,
		tlsCertFilePath:      tlsCertFilePath,
		tlsKeyFilePath:       tlsKeyFilePath,
	}
}

func (a *MainApplication) Run() {
	log.SetFlags(0)
	http.HandleFunc("/", a.SignallingController.WebSocket)
	log.Printf("Running on: https://%s", a.serverAddress)
	log.Fatal(
		http.ListenAndServeTLS(
			string(a.serverAddress),
			string(a.tlsCertFilePath),
			string(a.tlsKeyFilePath),
			nil,
		),
	)
}
