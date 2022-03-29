package database

import (
	"log"
	"react-demo-server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func autoMigrate(db *gorm.DB) {

	// if (!db.Migrator().HasTable(&model.User{})) {
	db.AutoMigrate(&model.User{})
	log.Println("创建用户表成功")
	// }

	if (!db.Migrator().HasTable(&model.Menu{})) {
		db.AutoMigrate(&model.Menu{})
		log.Println("创建菜单表成功")
	}

	// if (!db.Migrator().HasTable(&model.Role{})) {
	db.AutoMigrate(&model.Role{})
	log.Println("创建角色表成功")
	// }
}

func Start() *gorm.DB {

	logger := logger.Default.LogMode(logger.Info)

	db := PostgresStart(logger)

	autoMigrate(db)

	return db
}
