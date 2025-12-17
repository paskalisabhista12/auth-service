package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

var DB *gorm.DB

func Connect(dsn string) error {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = database
	slog.Info("connected to Postgres")
	return nil
}
