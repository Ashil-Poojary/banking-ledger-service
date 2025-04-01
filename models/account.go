package models

import (
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
	AccountType   string    `gorm:"not null" json:"account_type" validate:"oneof=Savings Checking Business"`
	Balance       float64   `gorm:"not null;default:0" json:"balance"`
	Currency      string    `gorm:"not null" json:"currency" validate:"len=3,uppercase"`
	CreatedAt     time.Time `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null;default:current_timestamp on update current_timestamp" json:"updated_at"`
}

// BeforeCreate sets UUIDs before inserting records
func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}
