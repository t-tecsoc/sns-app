package db

import (
	"backend/graph/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&model.Todo{})
	fmt.Println("migrated")
	return db
}
