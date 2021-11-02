package users

import (
	"errors"
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/interfaces"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
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

func prepareResponse(
	user *interfaces.User,
	accounts []interfaces.ResponseAccount,
	withToken bool,
) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var response = map[string]interface{}{
		"message": "All is fine",
		"status":  200,
	}

	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}

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

		result := db.Where("username = ?", username).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User not found", "status": 404}
		}

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong credentials", "status": 400}
		}

		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, true)

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
		return map[string]interface{}{"message": "Invalid username, email or password", "status": 400}
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

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{
			ID:      account.ID,
			Name:    account.Name,
			Balance: int(account.Balance),
		}
		accounts = append(accounts, respAccount)

		var response = prepareResponse(user, accounts, true)

		return response
	}
}

func GetUser(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)
	if isValid {
		db := helpers.ConnectDB()
		user := &interfaces.User{}

		result := db.Where("id = ? ", id).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User not found", "status": 404}
		}
		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, false)
		return response
	} else {
		return map[string]interface{}{"message": "Invalid token", "status": 400}
	}
}
