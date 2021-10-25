package users

import (
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/interfaces"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr((err))

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var token = prepareToken(user)
	var response = map[string]interface{}{
		"message": "All is fine",
		"status":  200,
	}

	response["jwt"] = token
	response["data"] = responseUser

	return response
}

func Login(username string, pass string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{
				Value: username, Name: "username",
			},
			{
				Value: pass, Name: "password",
			},
		},
	)

	if !valid {
		return map[string]interface{}{"message": "Invalid username or password", "status": 400}
	} else {
		db := helpers.ConnectDB()
		user := &interfaces.User{}

		if db.Where("username = ?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found", "status": 404}
		}

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong credentials", "status": 400}
		}

		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		defer db.Close()

		var response = prepareResponse(user, accounts)

		return response
	}
}

func Register(username string, email string, pass string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{
				Value: username, Name: "username",
			},
			{
				Value: email, Name: "email",
			},
			{
				Value: pass, Name: "password",
			},
		},
	)

	if !valid {
		return map[string]interface{}{"message": "Invalud username, email or password", "status": 400}
	} else {
		db := helpers.ConnectDB()

		generatedPassword := helpers.HashAndSalt([]byte(pass))

		user := &interfaces.User{
			Username: username,
			Email:    email,
			Password: generatedPassword,
		}

		db.Create(&user)

		account := &interfaces.Account{
			Type:    "Daily Account",
			Name:    string(username + "'s" + " account"),
			Balance: 0,
			UserID:  user.ID,
		}

		db.Create(&account)

		defer db.Close()

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{
			ID:      account.ID,
			Name:    account.Name,
			Balance: int(account.Balance),
		}
		accounts = append(accounts, respAccount)

		var response = prepareResponse(user, accounts)

		return response
	}
}
