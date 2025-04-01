package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction represents a bank transaction stored in MongoDB
type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`                          // Unique MongoDB ObjectID
	AccountID     string             `bson:"account_id" json:"account_id"`                     // Reference to Account ID
	AccountNumber string             `bson:"account_number" json:"account_number"`             // Account number
	Amount        float64            `bson:"amount" json:"amount"`                             // Transaction amount
	Type          string             `bson:"type" json:"type"`                                 // "deposit" or "withdrawal"
	Status        string             `bson:"status" json:"status"`                             // "pending", "completed", or "failed"
	Reference     string             `bson:"reference,omitempty" json:"reference,omitempty"`   // Optional transaction reference ID
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`                     // Timestamp when created
	UpdatedAt     time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"` // Timestamp when updated
}
