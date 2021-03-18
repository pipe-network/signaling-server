// +build ignore

package main

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/nacl/box"
	"os"
)

func main() {
	var err error
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("public.key", []byte(hex.EncodeToString(publicKey[:])), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("private.key", []byte(hex.EncodeToString(privateKey[:])), 0644)
	if err != nil {
		panic(err)
	}
}
