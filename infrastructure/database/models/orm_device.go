package models

import "gorm.io/gorm"

type ORMDevice struct {
	gorm.Model

	Token string
	PublicKey string
}
