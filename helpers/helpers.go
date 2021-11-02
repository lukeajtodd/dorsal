package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"luketodd/dorsal/interfaces"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	dsn := "host=127.0.0.1 port=5432 user=dorsal dbname=dorsal password=password sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	HandleErr(err)
	return db
}

func Validation(values []interfaces.Validation) bool {
	validate := validator.New()

	for i := 0; i < len(values); i++ {
		switch values[i].Name {
		case "username":
			errs := validate.Var(values[i].Value, "required")
			fmt.Println(errs)

			if errs != nil {
				return false
			}
		case "email":
			errs := validate.Var(values[i].Value, "required,email")
			fmt.Println(errs)

			if errs != nil {
				return false
			}
		case "password":
			errs := validate.Var(values[i].Value, "required")
			fmt.Println(errs)

			if errs != nil {
				return false
			}
		}
	}

	return true
}

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					log.Println(err)

					resp := interfaces.ErrResponse{Message: "Internal Server Error"}
					json.NewEncoder(w).Encode(resp)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}

func ValidateToken(id string, jwtToken string) bool {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})
	HandleErr(err)
	var userId, _ = strconv.ParseFloat(id, 32)
	if token.Valid && tokenData["user_id"] == userId {
		return true
	} else {
		return false
	}
}
