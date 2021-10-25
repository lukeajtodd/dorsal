package helpers

import (
	"luketodd/dorsal/interfaces"

	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}

func ConnectDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=dorsal dbname=dorsal password=password sslmode=disable")
	HandleErr(err)
	return db
}

func Validation(values []interfaces.Validation) bool {
	validate := validator.New()

	for i := 0; i < len(values); i++ {
		switch values[i].Name {
		case "username":
			errs := validate.Var(values[i].Value, "required,alphanum")

			if errs != nil {
				return false
			}
		case "email":
			errs := validate.Var(values[i].Value, "required,email")

			if errs != nil {
				return false
			}
		case "password":
			errs := validate.Var(values[i].Value, "required")

			if errs != nil {
				return false
			}
		}
	}

	return true
}
