package migrations

import (
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func connectDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=dorsal dbname=dorsal password=password sslmode=disable")
	helpers.HandleErr(err)
	return db
}

func createAccounts(db *gorm.DB) {
	users := [2]models.User{
		{
			Username: "Greg",
			Email:    "greg@gmail.com",
		},
		{
			Username: "Abigail",
			Email:    "abigail@gmail.com",
		},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))

		user := models.User{
			Username: users[1].Username,
			Email:    users[i].Email,
			Password: generatedPassword,
		}

		db.Create(&user)

		account := models.Account{
			Type:    "Daily Account",
			Name:    string(users[i].Username + "'s" + " account"),
			Balance: uint(10000 * int(i+1)),
			UserID:  user.ID,
		}

		db.Create(&account)
	}
}

func Migrate() {
	db := connectDB()
	db.AutoMigrate(&models.User{}, &models.Account{})
	defer db.Close()

	createAccounts(db)
}
