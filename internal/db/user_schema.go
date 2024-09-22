package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DomainID uint   `gorm:"not null; index;"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
}
