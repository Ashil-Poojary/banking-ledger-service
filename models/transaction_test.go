package models

import (
	"fmt"
	"testing"
	"time"
)

// TestTransactionValidation tests the transaction validation logic
func TestTransactionValidation(t *testing.T) {
	tests := []struct {
		name        string
		transaction Transaction
		expectError bool
	}{
		// ✅ Valid Deposit
		{
			name: "Valid deposit",
			transaction: Transaction{
				Type:          "deposit",
				AccountNumber: "12345",
				Amount:        100.00,
				Currency:      "USD",
				CreatedAt:     time.Now(),
			},
			expectError: false,
		},
		// ❌ Invalid Deposit (No Account Number)
		{
			name: "Invalid deposit (missing account number)",
			transaction: Transaction{
				Type:     "deposit",
				Amount:   50.00,
				Currency: "USD",
			},
			expectError: true,
		},
		// ✅ Valid Withdrawal
		{
			name: "Valid withdrawal",
			transaction: Transaction{
				Type:          "withdrawal",
				AccountNumber: "67890",
				Amount:        200.00,
				Currency:      "EUR",
				CreatedAt:     time.Now(),
			},
			expectError: false,
		},
		// ❌ Invalid Withdrawal (Negative Amount)
		{
			name: "Invalid withdrawal (negative amount)",
			transaction: Transaction{
				Type:          "withdrawal",
				AccountNumber: "67890",
				Amount:        -50.00,
				Currency:      "USD",
			},
			expectError: true,
		},
		// ✅ Valid Transfer
		{
			name: "Valid transfer",
			transaction: Transaction{
				Type:               "transfer",
				SourceAccount:      "12345",
				DestinationAccount: "67890",
				Amount:             300.00,
				Currency:           "GBP",
				CreatedAt:          time.Now(),
			},
			expectError: false,
		},
		// ❌ Invalid Transfer (Same Source & Destination)
		{
			name: "Invalid transfer (same source & destination account)",
			transaction: Transaction{
				Type:               "transfer",
				SourceAccount:      "11111",
				DestinationAccount: "11111",
				Amount:             100.00,
				Currency:           "USD",
			},
			expectError: true,
		},
		// ❌ Invalid Type
		{
			name: "Invalid transaction type",
			transaction: Transaction{
				Type:          "exchange",
				AccountNumber: "12345",
				Amount:        50.00,
				Currency:      "USD",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("Running test:", tc.name) // Print log
			err := tc.transaction.Validate()
			if err != nil {
				t.Logf("Transaction validation failed: %v", err) // Log error
			} else {
				t.Logf("Transaction validation succeeded for: %+v", tc.transaction) // Log success
			}

			// Check if error status matches expectation
			if (err != nil) != tc.expectError {
				t.Errorf("Test case '%s' failed: expected error %v, got %v", tc.name, tc.expectError, err)
			}
		})
	}
}
