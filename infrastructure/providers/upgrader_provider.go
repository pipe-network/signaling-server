package providers

import (
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/domain/values"
	"net/http"
)

func ProvideUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{
			values.SaltyRTCSubprotocol,
		},
	}
}
