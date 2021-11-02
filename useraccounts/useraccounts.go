package useraccounts

import (
	"fmt"
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/interfaces"
)

func updateAccount(id uint, amount int) {
	db := helpers.ConnectDB()
	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}

	db.Where("id = ? ", id).First(&account)
	account.Balance = uint(amount)
	db.Save(&account)

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)

	defer db.Close()
	return responseAcc
}

func getAccount(id *uint) *interfaces.Account {
	db := helpers.ConnectDB()
	account := &interfaces.Account{}
	if db.Where("id = ? ", id).First(&account).RecordNotFound() {
		return nil
	}

	defer db.Close()
	return account
}

func Transaction(userId uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
	userIdString := fmt.Sprint(userId)
	isValid := helpers.ValidateToken(userIdString, jwt)

	if isValid {
		fromAccount := getAccount(&from)
		toAccount := getAccount(&to)

		if fromAccount == nil || toAccount == nil {
			return map[string]interface{}{"message": "Account not found", "status": 404}
		} else if fromAccount.UserID != userId {
			return map[string]interface{}{"message": "You are not the owner of the account", "status": 400}]
		} else if int(fromAccount.Balance) < amount {
			return map[string]interface{}{"message": "Insufficient balance", "status": 400}
		}

		updatedAccount := updateAccount(from, int(fromAccount.Balance) - amount)
		updateAccount(to, int(toAccount.Balance) + amount)

		transactions.CreateTransaction(from, to, amount)

		var response = map[string]interface{}{"message": "Transfer complete", "status": 200}
		response["data"] = updatedAccount
		return response
	} else {
		return map[string]interface{}{"message": "Invalid token", "status": 400}
	}
}
