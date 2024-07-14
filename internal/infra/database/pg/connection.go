package db

import (
	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Connect(log logger.Logger, DSN string) (*gorm.DB, error) {
	con, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	product := &model.Product{}

	if err = con.AutoMigrate(product, user); err != nil {
		log.Error("migrate: error attempt to exec migrates")
		return nil, err
	}

	return con, nil
}
