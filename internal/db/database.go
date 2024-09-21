package db

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func Connect() (*Database, error) {
	DB_PATH := os.Getenv("DB_PATH")
	if DB_PATH == "" {
		DB_PATH = "courier.db"
	}

	db, err := gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	db.AutoMigrate(&Domain{}, &Address{})
	return &Database{db}, nil
}

func (d *Database) Disconnect() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
