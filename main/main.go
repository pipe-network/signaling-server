package main

import "github.com/pipe-network/signaling-server"

func main() {
	mainApplication := signaling_server.InitializeMainApplication()
	mainApplication.Run()
}
