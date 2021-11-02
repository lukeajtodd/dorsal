package migrations

import (
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/interfaces"

	"gorm.io/gorm"
)

func createAccounts(db *gorm.DB) {
	users := [2]interfaces.User{
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

		user := interfaces.User{
			Username: users[1].Username,
			Email:    users[i].Email,
			Password: generatedPassword,
		}

		db.Create(&user)

		account := interfaces.Account{
			Type:    "Daily Account",
			Name:    string(users[i].Username + "'s" + " account"),
			Balance: uint(10000 * int(i+1)),
			UserID:  user.ID,
		}

		db.Create(&account)
	}
}

func MigrateTransactions() {
	Transactions := &interfaces.Transaction{}

	db := helpers.ConnectDB()
	db.AutoMigrate(&Transactions)
}

func Migrate() {
	db := helpers.ConnectDB()
	db.AutoMigrate(&interfaces.User{}, &interfaces.Account{})

	createAccounts(db)
}
