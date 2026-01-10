package models

import "time"

// Balance represents account balance information following BIAN Account Balance Management domain
type Balance struct {
	// Balance identification and type
	BalanceType BalanceType `json:"balanceType"`
	
	// Balance amount
	Amount Money `json:"amount"`
	
	// Timestamp information
	Timestamp time.Time `json:"timestamp"`
}