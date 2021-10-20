package helpers

import (
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/models"

	"github.com/jinzhu/gorm"
)

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
