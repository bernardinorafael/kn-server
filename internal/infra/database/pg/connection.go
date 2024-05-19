package db

import (
	"log/slog"

	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Connect(l *slog.Logger, DSN string) (*gorm.DB, error) {
	l.Info("connecting database...")
	con, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	l.Info("generating migrations...")
	err = con.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, err
	}

	l.Info("database connected successfully")

	return con, nil
}
