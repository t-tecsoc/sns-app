package openapi

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init() {
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST_NAME")
	dBName := os.Getenv("MYSQL_DATABASE")

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	port := os.Getenv("PORT")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dBName + "?charset=utf8mb4"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}

	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = Db.AutoMigrate(&Users{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")
}
