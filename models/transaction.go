package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID string             `bson:"account_id" json:"account_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Type      string             `bson:"type" json:"type"` // deposit or withdrawal
	CreatedAt string             `bson:"created_at" json:"created_at"`
}
