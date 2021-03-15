package providers

import (
	"flag"
	"github.com/pipe-network/signaling-server/application"
	"github.com/pipe-network/signaling-server/infrastructure/storages"
)

func ProvideServerAddress() application.ServerAddress {
	address := flag.String("address", "localhost:8080", "http service address")
	return application.ServerAddress(*address)
}

func ProvidePublicKeyPath() storages.PublicKeyPath {
	publicKeyFile := flag.String("public_key_file", "./public.key", "public key file path")
	return storages.PublicKeyPath(*publicKeyFile)
}

func ProvidePrivateKeyPath() storages.PrivateKeyPath {
	privateKeyFile := flag.String("private_key_file", "./private.key", "private key file path")
	return storages.PrivateKeyPath(*privateKeyFile)
}
