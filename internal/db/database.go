package db

import (
	"fmt"
	"os"

	"github.com/jarqvi/courier/internal/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var Instance *Database

func Connect() error {
	DB_PATH := os.Getenv("DB_PATH")
	if DB_PATH == "" {
		DB_PATH = "courier.db"
	}

	db, err := gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.AutoMigrate(&Domain{}, &Address{}, &User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}

	log.Logger.Info("database connected and schema migrated successfully")

	Instance = &Database{db}

	return nil
}

func (d *Database) Disconnect() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Logger.Info("database connection closed successfully")

	return nil
}
