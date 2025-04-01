package models

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	// Validate account fields before inserting
	if err := a.Validate(); err != nil {
		return err
	}

	// Generate unique account number if not set
	if a.AccountNumber == "" {
		a.AccountNumber = generateAccountNumber()
	}
	return nil
}

// Validate ensures that the account has valid data.
func (a *Account) Validate() error {
	// Ensure Owner Name is not empty
	if a.OwnerName == "" {
		return errors.New("owner name is required")
	}

	// Ensure Account Number is not empty
	if a.AccountNumber == "" {
		return errors.New("account number is required")
	}

	// Ensure Account Type is not empty
	if a.AccountType == "" {
		return errors.New("account type is required")
	}

	// Normalize account type (title case)
	caser := cases.Title(language.English)
	a.AccountType = caser.String(strings.ToLower(a.AccountType))

	// Validate Account Type
	if !AllowedAccountTypes[a.AccountType] {
		return fmt.Errorf("invalid account type: must be Savings, Checking, or Business (got '%s')", a.AccountType)
	}

	// Ensure Balance is not negative
	if a.Balance < 0 {
		return errors.New("balance cannot be negative")
	}

	// Validate currency format (ISO 4217, e.g., USD, EUR)
	currencyRegex := regexp.MustCompile(`^[A-Z]{3}$`)
	if !currencyRegex.MatchString(a.Currency) {
		return errors.New("invalid currency format; must be a 3-letter ISO code (e.g., USD, EUR)")
	}

	return nil
}

// generateAccountNumber creates a cryptographically secure 10-digit account number.
func generateAccountNumber() string {
	var num uint64
	binary.Read(rand.Reader, binary.LittleEndian, &num)
	return fmt.Sprintf("%010d", num%10000000000) // Ensure it's a 10-digit number
}
