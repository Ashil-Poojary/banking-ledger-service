package models

import "gorm.io/gorm"

// Account represents a bank account for GORM
type Account struct {
	gorm.Model
	ID      uint    `json:"id"`
	Balance float64 `json:"balance"`
}
