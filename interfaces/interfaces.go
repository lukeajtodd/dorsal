package interfaces

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

type ResponseAccount struct {
	ID      uint
	Name    string
	Balance int
}

type ResponseUser struct {
	ID       uint
	Username string
	Email    string
	Accounts []ResponseAccount
}

type Validation struct {
	Value string
	Name  string
}

type ErrResponse struct {
	Message string
	Status  interface{}
}

type Transaction struct {
	gorm.Model
	From   uint
	To     uint
	Amount int
}

type ResponseTransaction struct {
	ID     uint
	From   uint
	To     uint
	Amount int
}
