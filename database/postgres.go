package database

import (
	"log"

	"react-demo-server/config"
	"react-demo-server/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) {

	if (!db.Migrator().HasTable(&model.User{})) {
		db.AutoMigrate(&model.User{})
		log.Println("创建用户表成功")
	}

	if (!db.Migrator().HasTable(&model.Menu{})) {
		db.AutoMigrate(&model.Menu{})
		log.Println("创建菜单表成功")
	}
}

func PostgresStart() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.Config.DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	autoMigrate(db)
	return db
}
