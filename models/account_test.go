package models

import (
	"fmt"
	"testing"
	"time"
)

// TestAccountValidation tests the account creation logic
func TestAccountValidation(t *testing.T) {
	tests := []struct {
		name        string
		account     Account
		expectError bool
	}{
		// ✅ Valid Account Creation
		{
			name: "Valid account creation",
			account: Account{
				AccountNumber: "123456",
				OwnerName:     "John Doe",
				AccountType:   "Savings", // ✅ Set a valid account type
				Balance:       1000.00,
				Currency:      "USD",
				CreatedAt:     time.Now(),
			},
			expectError: false,
		},

		// ❌ Invalid Account (Missing Account Number)
		{
			name: "Invalid account (missing account number)",
			account: Account{
				OwnerName:   "Alice",
				AccountType: "Checking", // ✅ Ensure account type is set
				Balance:     500.00,
				Currency:    "USD",
			},
			expectError: true,
		},

		// ❌ Invalid Account (Negative Balance)
		{
			name: "Invalid account (negative balance)",
			account: Account{
				AccountNumber: "678901",
				OwnerName:     "Bob",
				AccountType:   "Business", // ✅ Set a valid account type
				Balance:       -100.00,
				Currency:      "EUR",
			},
			expectError: true,
		},

		// ❌ Invalid Account (Missing Owner Name)
		{
			name: "Invalid account (missing owner name)",
			account: Account{
				AccountNumber: "999999",
				AccountType:   "Savings", // ✅ Ensure account type is set
				Balance:       100.00,
				Currency:      "GBP",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("Running test:", tc.name) // Print log
			err := tc.account.Validate()
			if err != nil {
				t.Logf("Account validation failed: %v", err) // Log error
			} else {
				t.Logf("Account validation succeeded for: %+v", tc.account) // Log success
			}

			// Check if error status matches expectation
			if (err != nil) != tc.expectError {
				t.Errorf("Test case '%s' failed: expected error %v, got %v", tc.name, tc.expectError, err)
			}
		})
	}
}
