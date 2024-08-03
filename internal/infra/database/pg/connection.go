package database

import (
	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Connect(log logger.Logger, DSN string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	tables := []interface{}{
		&gormodel.Team{},
		&gormodel.User{},
		&gormodel.Product{},
	}

	if err = db.AutoMigrate(tables...); err != nil {
		log.Error("migrate: error attempt to exec migrates", "httperr", err.Error())
		return nil, err
	}

	return db, nil
}
