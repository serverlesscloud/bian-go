package models

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestMoney_NewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   decimal.Decimal
		currency string
		wantErr  bool
	}{
		{
			name:     "valid AUD amount",
			amount:   decimal.NewFromFloat(100.50),
			currency: "AUD",
			wantErr:  false,
		},
		{
			name:     "valid USD amount",
			amount:   decimal.NewFromFloat(250.75),
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "invalid currency",
			amount:   decimal.NewFromFloat(100.00),
			currency: "XXX",
			wantErr:  true,
		},
		{
			name:     "lowercase currency converted to uppercase",
			amount:   decimal.NewFromFloat(100.00),
			currency: "aud",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !money.Amount.Equal(tt.amount) {
					t.Errorf("NewMoney() amount = %v, want %v", money.Amount, tt.amount)
				}
				expectedCurrency := "AUD"
				if tt.currency == "USD" {
					expectedCurrency = "USD"
				}
				if money.Currency != expectedCurrency {
					t.Errorf("NewMoney() currency = %v, want %v", money.Currency, expectedCurrency)
				}
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	money1, _ := NewMoney(decimal.NewFromFloat(100.50), "AUD")
	money2, _ := NewMoney(decimal.NewFromFloat(50.25), "AUD")
	moneyUSD, _ := NewMoney(decimal.NewFromFloat(100.00), "USD")

	result, err := money1.Add(money2)
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}
	
	expected := decimal.NewFromFloat(150.75)
	if !result.Amount.Equal(expected) {
		t.Errorf("Add() amount = %v, want %v", result.Amount, expected)
	}

	// Test currency mismatch
	_, err = money1.Add(moneyUSD)
	if err == nil {
		t.Error("Add() should return error for currency mismatch")
	}
}

func TestMoney_Subtract(t *testing.T) {
	money1, _ := NewMoney(decimal.NewFromFloat(100.50), "AUD")
	money2, _ := NewMoney(decimal.NewFromFloat(50.25), "AUD")

	result, err := money1.Subtract(money2)
	if err != nil {
		t.Errorf("Subtract() error = %v", err)
	}
	
	expected := decimal.NewFromFloat(50.25)
	if !result.Amount.Equal(expected) {
		t.Errorf("Subtract() amount = %v, want %v", result.Amount, expected)
	}
}

func TestMoney_Multiply(t *testing.T) {
	money, _ := NewMoney(decimal.NewFromFloat(100.00), "AUD")
	factor := decimal.NewFromFloat(1.5)

	result := money.Multiply(factor)
	expected := decimal.NewFromFloat(150.00)
	
	if !result.Amount.Equal(expected) {
		t.Errorf("Multiply() amount = %v, want %v", result.Amount, expected)
	}
}

func TestMoney_Divide(t *testing.T) {
	money, _ := NewMoney(decimal.NewFromFloat(100.00), "AUD")
	divisor := decimal.NewFromFloat(2.0)

	result, err := money.Divide(divisor)
	if err != nil {
		t.Errorf("Divide() error = %v", err)
	}
	
	expected := decimal.NewFromFloat(50.00)
	if !result.Amount.Equal(expected) {
		t.Errorf("Divide() amount = %v, want %v", result.Amount, expected)
	}

	// Test division by zero
	_, err = money.Divide(decimal.Zero)
	if err == nil {
		t.Error("Divide() should return error for division by zero")
	}
}

func TestMoney_JSONSerialization(t *testing.T) {
	money, _ := NewMoney(decimal.NewFromFloat(123.45), "AUD")

	// Test marshaling
	data, err := json.Marshal(money)
	if err != nil {
		t.Errorf("JSON Marshal error = %v", err)
	}

	expected := `{"amount":"123.45","currency":"AUD"}`
	if string(data) != expected {
		t.Errorf("JSON Marshal = %s, want %s", string(data), expected)
	}

	// Test unmarshaling
	var unmarshaled Money
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("JSON Unmarshal error = %v", err)
	}

	if !unmarshaled.Amount.Equal(money.Amount) || unmarshaled.Currency != money.Currency {
		t.Errorf("JSON Unmarshal = %+v, want %+v", unmarshaled, money)
	}
}

func TestMoney_Comparison(t *testing.T) {
	money1, _ := NewMoney(decimal.NewFromFloat(100.00), "AUD")
	money2, _ := NewMoney(decimal.NewFromFloat(100.00), "AUD")
	money3, _ := NewMoney(decimal.NewFromFloat(150.00), "AUD")

	// Test Equal
	if !money1.Equal(money2) {
		t.Error("Equal() should return true for same amounts")
	}

	// Test GreaterThan
	greater, err := money3.GreaterThan(money1)
	if err != nil {
		t.Errorf("GreaterThan() error = %v", err)
	}
	if !greater {
		t.Error("GreaterThan() should return true")
	}

	// Test LessThan
	less, err := money1.LessThan(money3)
	if err != nil {
		t.Errorf("LessThan() error = %v", err)
	}
	if !less {
		t.Error("LessThan() should return true")
	}
}