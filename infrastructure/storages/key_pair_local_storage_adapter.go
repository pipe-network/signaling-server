package storages

import (
	"encoding/hex"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/domain/values"
	"os"
	"strings"
)

type KeyPairLocalStorageAdapter struct {
	publicKey  values.Key
	privateKey values.Key

	publicKeyPath  string
	privateKeyPath string
}

func NewKeyPairLocalStorageAdapter(
	flagService services.FlagService,
) (*KeyPairLocalStorageAdapter, error) {

	keyPairStorageAdapter := &KeyPairLocalStorageAdapter{
		publicKeyPath:  flagService.String(services.PublicKeyFile),
		privateKeyPath: flagService.String(services.PrivateKeyFile),
	}

	err := keyPairStorageAdapter.Load()
	if err != nil {
		return nil, err
	}

	return keyPairStorageAdapter, nil
}

func (k *KeyPairLocalStorageAdapter) Load() error {
	privateKeyBytes, err := os.ReadFile(k.privateKeyPath)
	if err != nil {
		return err
	}
	privateKey := strings.TrimSpace(string(privateKeyBytes))
	decodedPrivateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	copy(k.privateKey[:], decodedPrivateKeyBytes[:])

	publicKeyBytes, err := os.ReadFile(k.publicKeyPath)
	if err != nil {
		return err
	}
	publicKey := strings.TrimSpace(string(publicKeyBytes))
	decodedPublicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return err
	}

	copy(k.publicKey[:], decodedPublicKeyBytes[:])
	return nil
}

func (k *KeyPairLocalStorageAdapter) PublicKey() values.Key {
	return k.publicKey
}

func (k *KeyPairLocalStorageAdapter) PrivateKey() values.Key {
	return k.privateKey
}
