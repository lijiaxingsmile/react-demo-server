package database

import (
	"react-demo-server/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func PostgresStart(log logger.Interface) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.Config.DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: log,
	})

	if err != nil {
		panic(err)
	}
	return db
}
