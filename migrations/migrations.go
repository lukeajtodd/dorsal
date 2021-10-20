package migrations

import (
	"luketodd/dorsal/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func connectDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=doral dbname=dorsal password=password sslmode=disable")
	helpers.HandleErr(err)
	return db
}
