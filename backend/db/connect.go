package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM(database_url string) *gorm.DB {

	dsn := os.Getenv(database_url)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	return db
}
