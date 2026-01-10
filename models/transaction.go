package models

import "time"

// Transaction represents a payment transaction following BIAN Payment Execution domain
type Transaction struct {
	// Transaction identification
	ID        string `json:"id"`
	Reference string `json:"reference,omitempty"`
	
	// Transaction classification
	TransactionType TransactionType `json:"transactionType"`
	
	// Transaction amount and currency
	Amount Money `json:"amount"`
	
	// Transaction details
	Description  string `json:"description"`
	MerchantName string `json:"merchantName,omitempty"`
	
	// Transaction dates
	PostingDate time.Time `json:"postingDate"`
	ValueDate   time.Time `json:"valueDate"`
	
	// Running balance after this transaction (optional)
	RunningBalance *Money `json:"runningBalance,omitempty"`
	
	// Account reference
	AccountID string `json:"accountId"`
}