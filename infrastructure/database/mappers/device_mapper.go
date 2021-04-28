package mappers

import (
	"github.com/pipe-network/signaling-server/domain/values"
	"github.com/pipe-network/signaling-server/infrastructure/database/models"
)

func MapDeviceToORMDevice(device values.Device) models.ORMDevice {
	return models.ORMDevice{
		Token:     device.Token,
		PublicKey: device.PublicKey,
	}
}

func MapORMDeviceToDevice(device models.ORMDevice) values.Device {
	return values.Device{
		Token:     device.Token,
		PublicKey: device.PublicKey,
	}
}
