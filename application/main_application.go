package application

import (
	"flag"
	"fmt"
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
	address := flag.String("address", "localhost", "http service host")
	port := flag.Int("port", 8080, "http service port")
	tlsCertFilePath := flag.String("tls_cert_file", "./cert.crt", "TLS certificate file path")
	tlsKeyFilePath := flag.String("tls_key_file", "./cert.key", "TLS key file path")

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", a.SignallingController.WebSocket)
	log.Printf("Running on: https://%s:%d", *address, *port)
	log.Fatal(
		http.ListenAndServeTLS(
			fmt.Sprintf("%s:%d", *address, *port),
			*tlsCertFilePath,
			*tlsKeyFilePath,
			nil,
		),
	)
}
