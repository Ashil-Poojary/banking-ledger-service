package models

import "time"

// Account represents a user account in the system
type Account struct {
	ID        int       `db:"id" json:"id"`
	Owner     string    `db:"owner" json:"owner"`
	Balance   float64   `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Transaction represents a financial transaction (deposit/withdrawal)
type Transaction struct {
	ID        int       `db:"id" json:"id"`
	AccountID int       `db:"account_id" json:"account_id"`
	Amount    float64   `db:"amount" json:"amount"`
	Type      string    `db:"type" json:"type"` // "deposit" or "withdrawal"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
