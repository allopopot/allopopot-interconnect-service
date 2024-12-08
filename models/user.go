package models

import (
	"allopopot-interconnect-service/service/passwordservice"
	"log"

	"gorm.io/gorm"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"-"`
	gorm.Model
}

func (u *User) SetPassword(password string) {
	hashedPassword, err := passwordservice.HashPassword(password)
	if err != nil {
		log.Panicln("Hashing Password Failed")
	}
	u.Password = hashedPassword
}

func (u *User) VerifyPassword(password string) bool {
	ok, _ := passwordservice.VerifyPassword(password, u.Password)
	return ok
}
