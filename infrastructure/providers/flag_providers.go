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

func ProvideTLSCertFilePath() application.TLSCertFilePath {
	tlsCertFilePath := flag.String("tls_cert_file", "./cert.pem", "TLS certificate file path")
	return application.TLSCertFilePath(*tlsCertFilePath)
}

func ProvideTLSKeyFilePath() application.TLSKeyFilePath {
	tlsKeyFilePath := flag.String("tls_key_file", "./key.pem", "TLS key file path")
	return application.TLSKeyFilePath(*tlsKeyFilePath)
}
