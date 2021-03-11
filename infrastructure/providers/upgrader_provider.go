package providers

import "github.com/gorilla/websocket"

func ProvideUpgrader() websocket.Upgrader {
	return websocket.Upgrader{}
}
