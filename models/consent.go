package models

import "time"

// Consent represents customer consent following BIAN Customer Consent Management domain
type Consent struct {
	// Consent identification
	ID string `json:"id"`
	
	// Consent status
	Status ConsentStatus `json:"status"`
	
	// Consent scopes (permissions granted)
	Scopes []string `json:"scopes"`
	
	// Consent lifecycle dates
	GrantDate      time.Time  `json:"grantDate"`
	ExpiryDate     time.Time  `json:"expiryDate"`
	RevocationDate *time.Time `json:"revocationDate,omitempty"`
}