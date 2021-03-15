package ports

import "github.com/pipe-network/signaling-server/domain/values"

type KeyPairStoragePort interface {
	Load() error
	PublicKey() values.Key
	PrivateKey() values.Key
}
