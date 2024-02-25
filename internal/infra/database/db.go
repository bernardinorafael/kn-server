package database

import (
	"log/slog"

	config "github.com/bernardinorafael/gozinho/config/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	env := config.Env

	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	slog.Info("Database connected!")

	return db, err
}
