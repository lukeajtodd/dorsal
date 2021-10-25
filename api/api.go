package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"luketodd/dorsal/helpers"
	"luketodd/dorsal/users"
	"net/http"

	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
}

type ErrResponse struct {
	Message string
	Status  interface{}
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	login := users.Login(formattedBody.Username, formattedBody.Password)

	if login["status"] == 200 {
		resp := login
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := ErrResponse{
			Message: "Wrong username or password",
			Status:  login["status"],
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	var formattedBody Register
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)

	if register["status"] == 200 {
		resp := register
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := ErrResponse{
			Message: "Invalid credentials",
			Status:  register["status"],
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func StartApi() {
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	fmt.Println("App is working on port :8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
