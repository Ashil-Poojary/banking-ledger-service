package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Allowed transaction types
var validTransactionTypes = map[string]bool{
	"deposit":    true,
	"withdrawal": true,
	"transfer":   true,
}

// Allowed currencies
var validCurrencies = map[string]bool{
	"USD": true, "EUR": true, "GBP": true, "INR": true, "JPY": true,
}

type TransferRequest struct {
	SourceAccount      string  `json:"source_account"`
	DestinationAccount string  `json:"destination_account"`
	Amount             float64 `json:"amount"`
	Currency           string  `json:"currency"`
}

// Transaction represents a bank transaction stored in MongoDB
type Transaction struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	SourceAccount      string             `bson:"source_account"`
	DestinationAccount string             `bson:"destination_account"`
	AccountNumber      string             `bson:"account_number"`
	Amount             float64            `bson:"amount"`
	Currency           string             `bson:"currency"`
	Type               string             `bson:"type"`
	Status             string             `bson:"status"`
	Reference          string             `bson:"reference"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
}

// Validate checks if the transaction data is valid
func (t *Transaction) Validate() error {
	// Validate transaction type
	if !validTransactionTypes[t.Type] {
		return errors.New("invalid transaction type: must be 'deposit', 'withdrawal', or 'transfer'")
	}

	// Validate amount (must be positive)
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Validate currency
	if t.Currency != "" && !validCurrencies[t.Currency] {
		return errors.New("invalid currency: must be a supported currency (e.g., USD, EUR, GBP)")
	}

	// Validate account fields based on transaction type
	switch t.Type {
	case "deposit", "withdrawal":
		if t.AccountNumber == "" {
			return errors.New("account number is required for deposits and withdrawals")
		}
		if t.SourceAccount != "" || t.DestinationAccount != "" {
			return errors.New("source and destination accounts should be empty for deposits and withdrawals")
		}

	case "transfer":
		if t.SourceAccount == "" || t.DestinationAccount == "" {
			return errors.New("source and destination accounts are required for transfers")
		}
		if t.AccountNumber != "" {
			return errors.New("account number should be empty for transfers")
		}
		if t.SourceAccount == t.DestinationAccount {
			return errors.New("source and destination accounts must be different")
		}
	}

	return nil
}
