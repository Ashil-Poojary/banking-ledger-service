package models

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account represents a user's bank account.
type Account struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE;" json:"user_id"` // Foreign key reference
	OwnerName     string    `gorm:"not null" json:"owner_name"`
	AccountNumber string    `gorm:"unique;not null" json:"account_number"`
	AccountType   string    `gorm:"not null" json:"account_type"`
	Balance       float64   `gorm:"not null;default:0" json:"balance"`
	Currency      string    `gorm:"not null" json:"currency"`
	CreatedAt     time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"not null;default:current_timestamp"`
}

// AllowedAccountTypes defines the valid types for an account.
var AllowedAccountTypes = map[string]bool{
	"Savings":  true,
	"Checking": true,
	"Business": true,
}

// BeforeCreate runs before inserting a new record.
func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	// Ensure UUID is generated
	a.ID = uuid.New()

	// Validate Account Type
	if !AllowedAccountTypes[a.AccountType] {
		return errors.New("invalid account type; must be Savings, Checking, or Business")
	}

	// Generate unique account number
	a.AccountNumber = generateAccountNumber()
	return nil
}

// generateAccountNumber creates a random 10-digit account number.
func generateAccountNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%010d", r.Intn(1000000000))
}
