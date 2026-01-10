package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/serverlesscloud/bian-go/domains"
	"github.com/serverlesscloud/bian-go/models"
	"github.com/shopspring/decimal"
)

// Provider implements all domain service interfaces with in-memory mock data
type Provider struct {
	accounts     map[string]*models.Account
	transactions map[string]*models.Transaction
	balances     map[string][]*models.Balance
	consents     map[string]*models.Consent
}

// NewProvider creates a new mock provider with sample data
func NewProvider() *Provider {
	p := &Provider{
		accounts:     make(map[string]*models.Account),
		transactions: make(map[string]*models.Transaction),
		balances:     make(map[string][]*models.Balance),
		consents:     make(map[string]*models.Consent),
	}
	p.loadSampleData()
	return p
}

// Ensure Provider implements all domain interfaces
var _ domains.AccountService = (*Provider)(nil)
var _ domains.TransactionService = (*Provider)(nil)
var _ domains.BalanceService = (*Provider)(nil)
var _ domains.ConsentService = (*Provider)(nil)

// AccountService implementation
func (p *Provider) RetrieveCurrentAccount(ctx context.Context, accountID string) (*models.Account, error) {
	account, exists := p.accounts[accountID]
	if !exists {
		return nil, fmt.Errorf("account not found: %s", accountID)
	}
	return account, nil
}

func (p *Provider) RetrieveCurrentAccountBalance(ctx context.Context, accountID string) (*models.Balance, error) {
	balances, exists := p.balances[accountID]
	if !exists {
		return nil, fmt.Errorf("account not found: %s", accountID)
	}
	
	// Return current balance
	for _, balance := range balances {
		if balance.BalanceType == models.BalanceTypeCurrent {
			return balance, nil
		}
	}
	
	return nil, fmt.Errorf("current balance not found for account: %s", accountID)
}

// TransactionService implementation
func (p *Provider) RetrievePaymentTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	transaction, exists := p.transactions[transactionID]
	if !exists {
		return nil, fmt.Errorf("transaction not found: %s", transactionID)
	}
	return transaction, nil
}

func (p *Provider) RetrievePaymentTransactionHistory(ctx context.Context, accountID string, opts domains.HistoryOptions) ([]*models.Transaction, error) {
	// Check if account exists
	if _, exists := p.accounts[accountID]; !exists {
		return nil, fmt.Errorf("account not found: %s", accountID)
	}
	
	var transactions []*models.Transaction
	
	// Filter transactions by account ID
	for _, tx := range p.transactions {
		if tx.AccountID == accountID {
			// Apply date filtering
			if opts.FromDate != nil && tx.PostingDate.Before(*opts.FromDate) {
				continue
			}
			if opts.ToDate != nil && tx.PostingDate.After(*opts.ToDate) {
				continue
			}
			transactions = append(transactions, tx)
		}
	}
	
	// Sort by posting date (newest first)
	for i := 0; i < len(transactions)-1; i++ {
		for j := i + 1; j < len(transactions); j++ {
			if transactions[i].PostingDate.Before(transactions[j].PostingDate) {
				transactions[i], transactions[j] = transactions[j], transactions[i]
			}
		}
	}
	
	// Apply pagination
	if opts.Offset > 0 {
		if opts.Offset >= len(transactions) {
			return []*models.Transaction{}, nil
		}
		transactions = transactions[opts.Offset:]
	}
	
	if opts.Limit > 0 && opts.Limit < len(transactions) {
		transactions = transactions[:opts.Limit]
	}
	
	return transactions, nil
}

// BalanceService implementation
func (p *Provider) RetrieveAccountBalance(ctx context.Context, accountID string) ([]*models.Balance, error) {
	balances, exists := p.balances[accountID]
	if !exists {
		return nil, fmt.Errorf("account not found: %s", accountID)
	}
	return balances, nil
}

// ConsentService implementation
func (p *Provider) RetrieveConsent(ctx context.Context, consentID string) (*models.Consent, error) {
	consent, exists := p.consents[consentID]
	if !exists {
		return nil, fmt.Errorf("consent not found: %s", consentID)
	}
	return consent, nil
}

func (p *Provider) RetrieveConsentStatus(ctx context.Context, consentID string) (models.ConsentStatus, error) {
	consent, exists := p.consents[consentID]
	if !exists {
		return "", fmt.Errorf("consent not found: %s", consentID)
	}
	return consent.Status, nil
}

// loadSampleData populates the provider with realistic test data
func (p *Provider) loadSampleData() {
	now := time.Now()
	
	// Sample accounts
	p.accounts["acc-001"] = &models.Account{
		ID:            "acc-001",
		AccountNumber: "123456789",
		AccountType:   models.AccountTypeChecking,
		ProductName:   "Everyday Checking",
		Nickname:      "Main Account",
		Status:        models.AccountStatusOpen,
		OpenDate:      now.AddDate(-2, 0, 0),
		Currency:      "AUD",
	}
	
	p.accounts["acc-002"] = &models.Account{
		ID:            "acc-002",
		AccountNumber: "987654321",
		AccountType:   models.AccountTypeSavings,
		ProductName:   "High Interest Savings",
		Status:        models.AccountStatusOpen,
		OpenDate:      now.AddDate(-1, -6, 0),
		Currency:      "AUD",
	}
	
	p.accounts["acc-003"] = &models.Account{
		ID:            "acc-003",
		AccountNumber: "555666777",
		AccountType:   models.AccountTypeCreditCard,
		ProductName:   "Platinum Credit Card",
		Status:        models.AccountStatusOpen,
		OpenDate:      now.AddDate(-3, 0, 0),
		Currency:      "USD",
	}
	
	// Sample balances
	currentBalance1, _ := models.NewMoney(decimal.NewFromFloat(2547.83), "AUD")
	availableBalance1, _ := models.NewMoney(decimal.NewFromFloat(2547.83), "AUD")
	
	p.balances["acc-001"] = []*models.Balance{
		{
			BalanceType: models.BalanceTypeCurrent,
			Amount:      *currentBalance1,
			Timestamp:   now,
		},
		{
			BalanceType: models.BalanceTypeAvailable,
			Amount:      *availableBalance1,
			Timestamp:   now,
		},
	}
	
	currentBalance2, _ := models.NewMoney(decimal.NewFromFloat(15420.50), "AUD")
	availableBalance2, _ := models.NewMoney(decimal.NewFromFloat(15420.50), "AUD")
	
	p.balances["acc-002"] = []*models.Balance{
		{
			BalanceType: models.BalanceTypeCurrent,
			Amount:      *currentBalance2,
			Timestamp:   now,
		},
		{
			BalanceType: models.BalanceTypeAvailable,
			Amount:      *availableBalance2,
			Timestamp:   now,
		},
	}
	
	currentBalance3, _ := models.NewMoney(decimal.NewFromFloat(-1250.75), "USD")
	availableBalance3, _ := models.NewMoney(decimal.NewFromFloat(3749.25), "USD")
	
	p.balances["acc-003"] = []*models.Balance{
		{
			BalanceType: models.BalanceTypeCurrent,
			Amount:      *currentBalance3,
			Timestamp:   now,
		},
		{
			BalanceType: models.BalanceTypeAvailable,
			Amount:      *availableBalance3,
			Timestamp:   now,
		},
	}
	
	// Sample transactions
	p.loadSampleTransactions(now)
	
	// Sample consents
	p.consents["consent-001"] = &models.Consent{
		ID:         "consent-001",
		Status:     models.ConsentStatusActive,
		Scopes:     []string{"account:read", "transaction:read", "balance:read"},
		GrantDate:  now.AddDate(0, -1, 0),
		ExpiryDate: now.AddDate(0, 11, 0),
	}
	
	expiredDate := now.AddDate(0, -2, 0)
	p.consents["consent-002"] = &models.Consent{
		ID:         "consent-002",
		Status:     models.ConsentStatusExpired,
		Scopes:     []string{"account:read", "balance:read"},
		GrantDate:  now.AddDate(0, -14, 0),
		ExpiryDate: expiredDate,
	}
	
	revokedDate := now.AddDate(0, 0, -15)
	p.consents["consent-003"] = &models.Consent{
		ID:             "consent-003",
		Status:         models.ConsentStatusRevoked,
		Scopes:         []string{"account:read", "transaction:read"},
		GrantDate:      now.AddDate(0, -3, 0),
		ExpiryDate:     now.AddDate(0, 9, 0),
		RevocationDate: &revokedDate,
	}
}

func (p *Provider) loadSampleTransactions(now time.Time) {
	transactions := []struct {
		id          string
		accountID   string
		txType      models.TransactionType
		amount      string
		currency    string
		description string
		merchant    string
		daysAgo     int
	}{
		{"tx-001", "acc-001", models.TransactionTypeDebit, "-45.50", "AUD", "Grocery Shopping", "Woolworths", 1},
		{"tx-002", "acc-001", models.TransactionTypeCredit, "2500.00", "AUD", "Salary Payment", "ACME Corp", 3},
		{"tx-003", "acc-001", models.TransactionTypeDebit, "-12.80", "AUD", "Coffee", "Local Cafe", 5},
		{"tx-004", "acc-001", models.TransactionTypeDebit, "-89.99", "AUD", "Utilities", "Energy Australia", 7},
		{"tx-005", "acc-001", models.TransactionTypeDebit, "-25.00", "AUD", "ATM Withdrawal", "ANZ ATM", 10},
		
		{"tx-006", "acc-002", models.TransactionTypeCredit, "500.00", "AUD", "Transfer from Checking", "Internal Transfer", 2},
		{"tx-007", "acc-002", models.TransactionTypeCredit, "15.25", "AUD", "Interest Payment", "Bank Interest", 30},
		{"tx-008", "acc-002", models.TransactionTypeCredit, "1000.00", "AUD", "Bonus Payment", "ACME Corp", 45},
		
		{"tx-009", "acc-003", models.TransactionTypeDebit, "-125.00", "USD", "Online Shopping", "Amazon", 2},
		{"tx-010", "acc-003", models.TransactionTypeDebit, "-89.50", "USD", "Restaurant", "Fine Dining Co", 5},
		{"tx-011", "acc-003", models.TransactionTypeCredit, "-25.00", "USD", "Payment Received", "Credit Payment", 15},
	}
	
	for _, tx := range transactions {
		amount, _ := models.NewMoneyFromString(tx.amount, tx.currency)
		postingDate := now.AddDate(0, 0, -tx.daysAgo)
		
		p.transactions[tx.id] = &models.Transaction{
			ID:              tx.id,
			Reference:       fmt.Sprintf("REF-%s", tx.id),
			TransactionType: tx.txType,
			Amount:          *amount,
			Description:     tx.description,
			MerchantName:    tx.merchant,
			PostingDate:     postingDate,
			ValueDate:       postingDate,
			AccountID:       tx.accountID,
		}
	}
}