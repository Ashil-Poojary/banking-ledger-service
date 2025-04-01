package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a system user.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Phone     string    `gorm:"not null" json:"phone"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp on update current_timestamp" json:"updated_at"`
}

// BeforeCreate hashes the password and generates a UUID before inserting the user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

// CheckPassword compares the hashed password with the provided one
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
