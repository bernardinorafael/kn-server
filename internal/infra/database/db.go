package database

import (
	"log/slog"

	config "github.com/bernardinorafael/gozinho/config/env"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	env := config.Env

	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		slog.Error("Error connecting database!", err)
		return nil, err
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		slog.Error("Error generate migrations", err)
		return nil, err
	}

	slog.Info("Database connected!")

	return db, err
}
