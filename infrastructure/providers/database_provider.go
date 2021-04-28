package providers

import (
	"github.com/pipe-network/signaling-server/infrastructure/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DatabaseProvider() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.ORMDevice{})
	if err != nil {
		panic(err)
	}
	return db
}
