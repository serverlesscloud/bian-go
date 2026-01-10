package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Money represents a monetary amount with currency using decimal arithmetic
// to avoid floating-point precision errors in financial calculations.
type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// NewMoney creates a new Money instance with validation
func NewMoney(amount decimal.Decimal, currency string) (*Money, error) {
	if err := validateCurrency(currency); err != nil {
		return nil, err
	}
	return &Money{
		Amount:   amount,
		Currency: strings.ToUpper(currency),
	}, nil
}

// NewMoneyFromString creates Money from string amount and currency
func NewMoneyFromString(amount, currency string) (*Money, error) {
	dec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}
	return NewMoney(dec, currency)
}

// Add returns a new Money instance with the sum of two amounts
func (m *Money) Add(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, fmt.Errorf("currency mismatch: %s != %s", m.Currency, other.Currency)
	}
	return &Money{
		Amount:   m.Amount.Add(other.Amount),
		Currency: m.Currency,
	}, nil
}

// Subtract returns a new Money instance with the difference
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, fmt.Errorf("currency mismatch: %s != %s", m.Currency, other.Currency)
	}
	return &Money{
		Amount:   m.Amount.Sub(other.Amount),
		Currency: m.Currency,
	}, nil
}

// Multiply returns a new Money instance multiplied by a decimal factor
func (m *Money) Multiply(factor decimal.Decimal) *Money {
	return &Money{
		Amount:   m.Amount.Mul(factor),
		Currency: m.Currency,
	}
}

// Divide returns a new Money instance divided by a decimal divisor
func (m *Money) Divide(divisor decimal.Decimal) (*Money, error) {
	if divisor.IsZero() {
		return nil, fmt.Errorf("division by zero")
	}
	return &Money{
		Amount:   m.Amount.Div(divisor),
		Currency: m.Currency,
	}, nil
}

// Equal checks if two Money instances are equal
func (m *Money) Equal(other *Money) bool {
	return m.Currency == other.Currency && m.Amount.Equal(other.Amount)
}

// GreaterThan checks if this Money is greater than another
func (m *Money) GreaterThan(other *Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, fmt.Errorf("currency mismatch: %s != %s", m.Currency, other.Currency)
	}
	return m.Amount.GreaterThan(other.Amount), nil
}

// LessThan checks if this Money is less than another
func (m *Money) LessThan(other *Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, fmt.Errorf("currency mismatch: %s != %s", m.Currency, other.Currency)
	}
	return m.Amount.LessThan(other.Amount), nil
}

// IsZero checks if the amount is zero
func (m *Money) IsZero() bool {
	return m.Amount.IsZero()
}

// IsNegative checks if the amount is negative
func (m *Money) IsNegative() bool {
	return m.Amount.IsNegative()
}

// String returns a human-readable representation
func (m *Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount.String(), m.Currency)
}

// MarshalJSON implements json.Marshaler to ensure precision
func (m *Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}{
		Amount:   m.Amount.String(),
		Currency: m.Currency,
	})
}

// UnmarshalJSON implements json.Unmarshaler
func (m *Money) UnmarshalJSON(data []byte) error {
	var temp struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	
	amount, err := decimal.NewFromString(temp.Amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}
	
	if err := validateCurrency(temp.Currency); err != nil {
		return err
	}
	
	m.Amount = amount
	m.Currency = strings.ToUpper(temp.Currency)
	return nil
}

// validateCurrency checks if the currency code is valid ISO 4217
func validateCurrency(currency string) error {
	currency = strings.ToUpper(currency)
	validCurrencies := map[string]bool{
		"AUD": true, "USD": true, "GBP": true, "EUR": true, "CAD": true,
		"JPY": true, "CHF": true, "CNY": true, "SEK": true, "NZD": true,
	}
	
	if !validCurrencies[currency] {
		return fmt.Errorf("invalid currency code: %s", currency)
	}
	return nil
}