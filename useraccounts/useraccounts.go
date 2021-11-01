package useraccounts

import (
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/interfaces"
)

func updateAccount(id uint, amount int) {
	db := helpers.ConnectDB()
	db.Model(&interfaces.Account{}).Where("id = ? ", id).Update("balance", amount)
	defer db.Close()
}
