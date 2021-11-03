package database

import (
	"luketodd/dorsal/helpers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := "host=127.0.0.1 port=5432 user=dorsal dbname=dorsal password=password sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	helpers.HandleErr(err)

	pqDB, err := db.DB()

	helpers.HandleErr(err)

	pqDB.SetMaxIdleConns(20)
	pqDB.SetMaxOpenConns(200)

	DB = db
}
