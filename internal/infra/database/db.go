package database

import (
	"context"

	"github.com/bernardinorafael/gozinho/config"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	utillog "github.com/bernardinorafael/gozinho/util/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, l utillog.Logger) (*gorm.DB, error) {
	env := config.Env

	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		l.Errorf(ctx, "error connecting database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&entity.Account{})
	if err != nil {
		l.Errorf(ctx, "error generate migrations: %v", err)
		return nil, err
	}

	l.Info(ctx, "database connected successfully")

	return db, err
}
