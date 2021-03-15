package main

import "github.com/pipe-network/signaling-server"

func main() {
	mainApplication, err := signaling_server.InitializeMainApplication()
	if err != nil {
		panic(err)
	}

	mainApplication.Run()
}
