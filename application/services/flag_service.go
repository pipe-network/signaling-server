package services

import (
	"flag"
)

const (
	Address        = "address"
	Port           = "port"
	TLSCertFile    = "tls_cert_file"
	TLSKeyFile     = "tls_key_file"
	FCMServerKey   = "fcm_server_key"
	PublicKeyFile  = "public_key_file"
	PrivateKeyFile = "private_key_file"
)

type (
	FlagService interface {
		String(key string) string
		Int(key string) int
	}
	FlagServiceImpl struct {
		stringFlags map[string]string
		intFlags    map[string]int
	}
)

func NewFlagServiceImpl() FlagService {
	flagServiceImpl := &FlagServiceImpl{}
	flagServiceImpl.init()
	return flagServiceImpl
}

func (i *FlagServiceImpl) init() {
	i.intFlags = make(map[string]int)
	i.stringFlags = make(map[string]string)

	address := flag.String(Address, "localhost", "http service host")
	tlsCertFilePath := flag.String(TLSCertFile, "./cert.crt", "TLS certificate file path")
	tlsKeyFilePath := flag.String(TLSKeyFile, "./cert.key", "TLS key file path")
	fcmServerKey := flag.String(FCMServerKey, "", "FCM Server key")
	publicKeyPath := flag.String(PublicKeyFile, "./public.key", "public key file path")
	privateKeyPath := flag.String(PrivateKeyFile, "./private.key", "private key file path")
	port := flag.Int(Port, 8080, "http service port")

	flag.Parse()

	i.stringFlags[Address] = *address
	i.stringFlags[TLSCertFile] = *tlsCertFilePath
	i.stringFlags[TLSKeyFile] = *tlsKeyFilePath
	i.stringFlags[FCMServerKey] = *fcmServerKey
	i.stringFlags[PublicKeyFile] = *publicKeyPath
	i.stringFlags[PrivateKeyFile] = *privateKeyPath
	i.intFlags[Port] = *port
}

func (i *FlagServiceImpl) String(key string) string {
	return i.stringFlags[key]
}

func (i *FlagServiceImpl) Int(key string) int {
	return i.intFlags[key]
}
