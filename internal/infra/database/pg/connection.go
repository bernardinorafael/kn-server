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
	var user user.User
	var product product.Product

	con, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = con.AutoMigrate(&user, &product)
	if err != nil {
		return nil, err
	}

	l.Info("database connected")
	return con, nil
}
