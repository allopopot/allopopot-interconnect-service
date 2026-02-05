package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpencierProjects struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `json:"name" bson:"name"`
	Description   string             `json:"description" bson:"description"`
	SubLists      []string           `json:"sub_lists" bson:"sub_lists"`
	UserId        primitive.ObjectID `json:"user_id" bson:"user_id"`
	CreatedTime   time.Time          `json:"created_time" bson:"created_time"`
	CurrentAmount float64            `json:"current_amount" bson:"current_amount"`
}
