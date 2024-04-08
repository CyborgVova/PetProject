package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Mapping struct {
	Short string
	Long  string
}

func CreateDB(storage string) *gorm.DB {
	if storage == "inmemory" {
		return nil
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("MY_HOST"),
		os.Getenv("MY_LOGIN"),
		os.Getenv("MY_PASS"),
		os.Getenv("MY_BASE"),
		os.Getenv("MY_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Open database ", err)
	}
	db.AutoMigrate(&Mapping{})
	return db
}
