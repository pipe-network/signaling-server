package providers

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	SaltyRTCSubprotocol = "v1.saltyrtc.org"
)

func ProvideUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{
			SaltyRTCSubprotocol,
		},
	}
}
