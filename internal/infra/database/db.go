package database

import (
	"context"

	"github.com/bernardinorafael/kn-server/config"
	utillog "github.com/bernardinorafael/kn-server/helper/log"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, log utillog.Logger) (*gorm.DB, error) {
	env := config.Env

	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		log.Errorf(ctx, "error connecting database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&entity.Account{})
	if err != nil {
		log.Errorf(ctx, "error generate migrations: %v", err)
		return nil, err
	}

	log.Info(ctx, "database connected")

	return db, err
}
