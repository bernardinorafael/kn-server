package db

import (
	"log/slog"

	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
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
	err = con.AutoMigrate(&user.User{}, &product.Product{})
	if err != nil {
		return nil, err
	}

	l.Info("database connected successfully")

	return con, nil
}
