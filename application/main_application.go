package application

import (
	"flag"
	"github.com/pipe-network/signaling-server/interface/controllers"
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
	address := flag.String("address", "localhost:8080", "http service address")
	tlsCertFilePath := flag.String("tls_cert_file", "./cert.pem", "TLS certificate file path")
	tlsKeyFilePath := flag.String("tls_key_file", "./key.pem", "TLS key file path")
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", a.SignallingController.WebSocket)
	log.Printf("Running on: https://%s", *address)
	log.Fatal(
		http.ListenAndServeTLS(
			*address,
			*tlsCertFilePath,
			*tlsKeyFilePath,
			nil,
		),
	)
}
