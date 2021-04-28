package repositories

import (
	"errors"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/domain/values"
	"github.com/pipe-network/signaling-server/infrastructure/database/mappers"
	"github.com/pipe-network/signaling-server/infrastructure/database/models"
	"gorm.io/gorm"
)

type DeviceTokenDatabaseRepository struct {
	database *gorm.DB
}

func NewDeviceTokenDatabaseRepository(database *gorm.DB) ports.DeviceTokenRepository {
	return &DeviceTokenDatabaseRepository{database: database}
}

func (d *DeviceTokenDatabaseRepository) CreateOrUpdateToken(device values.Device) error {
	ormDevice := &models.ORMDevice{}
	result := d.database.First(&ormDevice, "public_key = ?", device.PublicKey)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newOrmDevice := mappers.MapDeviceToORMDevice(device)
		result = d.database.Create(&newOrmDevice)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}

	ormDevice.Token = device.Token
	result = d.database.Save(&ormDevice)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DeviceTokenDatabaseRepository) DeviceByPublicKey(publicKeyHex string) (*values.Device, error) {
	ormDevice := models.ORMDevice{}
	result := d.database.First(&ormDevice, "public_key = ?", publicKeyHex)
	if result.Error != nil {
		return nil, result.Error
	}
	device := mappers.MapORMDeviceToDevice(ormDevice)
	return &device, nil
}
