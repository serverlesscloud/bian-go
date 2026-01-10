package models

import "time"

// Account represents a bank account following BIAN Current Account Fulfillment domain
type Account struct {
	// Account identification
	ID            string `json:"id"`
	AccountNumber string `json:"accountNumber"`
	
	// Account classification
	AccountType AccountType   `json:"accountType"`
	ProductName string        `json:"productName"`
	Nickname    string        `json:"nickname,omitempty"`
	
	// Account status and lifecycle
	Status    AccountStatus `json:"status"`
	OpenDate  time.Time     `json:"openDate"`
	CloseDate *time.Time    `json:"closeDate,omitempty"`
	
	// Additional metadata
	Currency string `json:"currency"`
}