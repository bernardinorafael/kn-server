package database

import (
	"context"
	"log/slog"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, l *slog.Logger) (*gorm.DB, error) {
	env := config.Env

	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		l.Error("error connecting database", err)
		return nil, err
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		l.Error("error generate migrations: %v", err)
		return nil, err
	}

	l.Info("database connected!")

	return db, err
}
