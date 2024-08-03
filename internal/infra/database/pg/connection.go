package database

import (
	"sync"

	"gorm.io/driver/postgres"

	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Connect(log logger.Logger, DSN string) (*gorm.DB, error) {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
		if err != nil {
			return
		}

		tables := []interface{}{
			&gormodel.Team{},
			&gormodel.User{},
			&gormodel.Product{},
		}

		if err = db.AutoMigrate(tables...); err != nil {
			log.Error("migrate: error attempt to exec migrates", "httperr", err.Error())
			return
		}
	})

	return db, nil
}
