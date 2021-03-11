package main

import (
	"github.com/pipe-network/signaling-server/src"
)

func main() {
	mainApplication := src.InitializeMainApplication()
	mainApplication.Run()
}
