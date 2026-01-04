package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpencierTransactions struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProjectId   primitive.ObjectID `json:"project_id" bson:"project_id"`
	SubList     string             `json:"sub_list" bson:"sub_list"`
	Amount      float64            `json:"amount" bson:"amount"`
	CreatedTime time.Time          `json:"created_time" bson:"created_time"`
	Description string             `json:"description" bson:"description"`
}
