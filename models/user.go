package models

import (
	"allopopot-interconnect-service/service/passwordservice"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"-" bson:"password"`
	RecoveryCode string             `json:"-" bson:"recovery_code"`
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

func (u *User) GenerateRecoveryCode() string {
	recoveryCode := uuid.NewString()
	u.RecoveryCode = recoveryCode
	return recoveryCode
}
