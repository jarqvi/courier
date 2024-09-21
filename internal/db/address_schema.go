package db

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Name     string `gorm:"not null"`
	DomainID uint   `gorm:"not null; index;"`
}
