package db

import (
	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Connect(log logger.Logger, DSN string) (*gorm.DB, error) {
	con, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var team = &gormodel.Team{}
	var user = &gormodel.User{}
	var product = &gormodel.Product{}

	tables := []interface{}{
		team,
		user,
		product,
	}

	if err = con.AutoMigrate(tables...); err != nil {
		log.Error("migrate: error attempt to exec migrates", "error", err.Error())
		return nil, err
	}

	return con, nil
}
