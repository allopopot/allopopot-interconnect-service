package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Type        string             `json:"type" bson:"type"`
	Token       string             `json:"token" bson:"token"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	CreatedTime time.Time          `json:"created_time" bson:"created_time"`
	ExpiryTime  time.Time          `json:"expiry_time" bson:"expiry_time"`
}

func (t *Token) StoreToken(t_type string, token string) {

}
