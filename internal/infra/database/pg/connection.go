package db

import (
	"log/slog"

	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Connect(l *slog.Logger, DSN string) (*gorm.DB, error) {
	con, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = con.AutoMigrate(&model.Product{}, &model.User{})
	if err != nil {
		return nil, err
	}

	l.Info("database connected")
	return con, nil
}
