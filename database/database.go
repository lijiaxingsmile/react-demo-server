package database

import "gorm.io/gorm"

func Start() *gorm.DB {
	return PostgresStart()
}
