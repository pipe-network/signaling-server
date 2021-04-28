package ports

import "github.com/pipe-network/signaling-server/domain/values"

type DeviceTokenRepository interface {
	CreateOrUpdateToken(device values.Device) error
	DeviceByPublicKey(publicKeyHex string) (*values.Device, error)
}
