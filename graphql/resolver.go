package graphql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/serverlesscloud/bian-go/domains"
	"github.com/serverlesscloud/bian-go/graphql/generated"
	"github.com/serverlesscloud/bian-go/models"
)

// Resolver contains the domain services
type Resolver struct {
	accountService     domains.AccountService
	transactionService domains.TransactionService
	balanceService     domains.BalanceService
	consentService     domains.ConsentService
}

// NewResolver creates a new GraphQL resolver
func NewResolver(
	accountService domains.AccountService,
	transactionService domains.TransactionService,
	balanceService domains.BalanceService,
	consentService domains.ConsentService,
) *Resolver {
	return &Resolver{
		accountService:     accountService,
		transactionService: transactionService,
		balanceService:     balanceService,
		consentService:     consentService,
	}
}

// Query resolver implementation
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

// Account resolves the account query
func (r *queryResolver) Account(ctx context.Context, id string) (*generated.Account, error) {
	account, err := r.accountService.RetrieveCurrentAccount(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("account not found: %s", id)
		}
		return nil, err
	}
	
	return &generated.Account{
		ID:            account.ID,
		AccountNumber: account.AccountNumber,
		AccountType:   mapAccountType(account.AccountType),
		ProductName:   account.ProductName,
		Nickname:      &account.Nickname,
		Status:        mapAccountStatus(account.Status),
		OpenDate:      account.OpenDate.Format(time.RFC3339),
		CloseDate:     formatTimePtr(account.CloseDate),
		Currency:      account.Currency,
	}, nil
}

// Balance resolves the balance query (current balance only)
func (r *queryResolver) Balance(ctx context.Context, accountID string) (*generated.Balance, error) {
	balance, err := r.accountService.RetrieveCurrentAccountBalance(ctx, accountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("account not found: %s", accountID)
		}
		return nil, err
	}
	
	return &generated.Balance{
		BalanceType: mapBalanceType(balance.BalanceType),
		Amount: &generated.Money{
			Amount:   balance.Amount.Amount.String(),
			Currency: balance.Amount.Currency,
		},
		Timestamp: balance.Timestamp.Format(time.RFC3339),
	}, nil
}

// Balances resolves the balances query (all balance types)
func (r *queryResolver) Balances(ctx context.Context, accountID string) ([]*generated.Balance, error) {
	balances, err := r.balanceService.RetrieveAccountBalance(ctx, accountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("account not found: %s", accountID)
		}
		return nil, err
	}
	
	var result []*generated.Balance
	for _, balance := range balances {
		result = append(result, &generated.Balance{
			BalanceType: mapBalanceType(balance.BalanceType),
			Amount: &generated.Money{
				Amount:   balance.Amount.Amount.String(),
				Currency: balance.Amount.Currency,
			},
			Timestamp: balance.Timestamp.Format(time.RFC3339),
		})
	}
	
	return result, nil
}

// Transaction resolves the transaction query
func (r *queryResolver) Transaction(ctx context.Context, id string) (*generated.Transaction, error) {
	transaction, err := r.transactionService.RetrievePaymentTransaction(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("transaction not found: %s", id)
		}
		return nil, err
	}
	
	return mapTransaction(transaction), nil
}

// Transactions resolves the transactions query
func (r *queryResolver) Transactions(ctx context.Context, accountID string, input *generated.TransactionHistoryInput) ([]*generated.Transaction, error) {
	opts := domains.HistoryOptions{}
	
	if input != nil {
		if input.FromDate != nil {
			fromDate, err := time.Parse("2006-01-02", *input.FromDate)
			if err != nil {
				return nil, fmt.Errorf("invalid fromDate format: %w", err)
			}
			opts.FromDate = &fromDate
		}
		
		if input.ToDate != nil {
			toDate, err := time.Parse("2006-01-02", *input.ToDate)
			if err != nil {
				return nil, fmt.Errorf("invalid toDate format: %w", err)
			}
			opts.ToDate = &toDate
		}
		
		if input.Limit != nil {
			if *input.Limit <= 0 || *input.Limit > 500 {
				return nil, fmt.Errorf("limit must be between 1 and 500")
			}
			opts.Limit = *input.Limit
		}
		
		if input.Offset != nil {
			if *input.Offset < 0 {
				return nil, fmt.Errorf("offset must be non-negative")
			}
			opts.Offset = *input.Offset
		}
	}
	
	transactions, err := r.transactionService.RetrievePaymentTransactionHistory(ctx, accountID, opts)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("account not found: %s", accountID)
		}
		return nil, err
	}
	
	var result []*generated.Transaction
	for _, tx := range transactions {
		result = append(result, mapTransaction(tx))
	}
	
	return result, nil
}

// Consent resolves the consent query
func (r *queryResolver) Consent(ctx context.Context, id string) (*generated.Consent, error) {
	consent, err := r.consentService.RetrieveConsent(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("consent not found: %s", id)
		}
		return nil, err
	}
	
	return &generated.Consent{
		ID:             consent.ID,
		Status:         mapConsentStatus(consent.Status),
		Scopes:         consent.Scopes,
		GrantDate:      consent.GrantDate.Format(time.RFC3339),
		ExpiryDate:     consent.ExpiryDate.Format(time.RFC3339),
		RevocationDate: formatTimePtr(consent.RevocationDate),
	}, nil
}

// ConsentStatus resolves the consentStatus query
func (r *queryResolver) ConsentStatus(ctx context.Context, id string) (*generated.ConsentStatus, error) {
	status, err := r.consentService.RetrieveConsentStatus(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("consent not found: %s", id)
		}
		return nil, err
	}
	
	mappedStatus := mapConsentStatus(status)
	return &mappedStatus, nil
}

// Helper functions for mapping between domain models and GraphQL types

func mapAccountType(at models.AccountType) generated.AccountType {
	switch at {
	case models.AccountTypeChecking:
		return generated.AccountTypeChecking
	case models.AccountTypeSavings:
		return generated.AccountTypeSavings
	case models.AccountTypeCreditCard:
		return generated.AccountTypeCreditCard
	case models.AccountTypeInvestment:
		return generated.AccountTypeInvestment
	default:
		return generated.AccountTypeChecking
	}
}

func mapAccountStatus(as models.AccountStatus) generated.AccountStatus {
	switch as {
	case models.AccountStatusOpen:
		return generated.AccountStatusOpen
	case models.AccountStatusClosed:
		return generated.AccountStatusClosed
	case models.AccountStatusSuspended:
		return generated.AccountStatusSuspended
	default:
		return generated.AccountStatusOpen
	}
}

func mapBalanceType(bt models.BalanceType) generated.BalanceType {
	switch bt {
	case models.BalanceTypeCurrent:
		return generated.BalanceTypeCurrent
	case models.BalanceTypeAvailable:
		return generated.BalanceTypeAvailable
	case models.BalanceTypePending:
		return generated.BalanceTypePending
	default:
		return generated.BalanceTypeCurrent
	}
}

func mapTransactionType(tt models.TransactionType) generated.TransactionType {
	switch tt {
	case models.TransactionTypeDebit:
		return generated.TransactionTypeDebit
	case models.TransactionTypeCredit:
		return generated.TransactionTypeCredit
	case models.TransactionTypeTransfer:
		return generated.TransactionTypeTransfer
	case models.TransactionTypePayment:
		return generated.TransactionTypePayment
	case models.TransactionTypeFee:
		return generated.TransactionTypeFee
	default:
		return generated.TransactionTypeDebit
	}
}

func mapConsentStatus(cs models.ConsentStatus) generated.ConsentStatus {
	switch cs {
	case models.ConsentStatusActive:
		return generated.ConsentStatusActive
	case models.ConsentStatusExpired:
		return generated.ConsentStatusExpired
	case models.ConsentStatusRevoked:
		return generated.ConsentStatusRevoked
	case models.ConsentStatusPending:
		return generated.ConsentStatusPending
	default:
		return generated.ConsentStatusActive
	}
}

func mapTransaction(tx *models.Transaction) *generated.Transaction {
	result := &generated.Transaction{
		ID:              tx.ID,
		Reference:       &tx.Reference,
		TransactionType: mapTransactionType(tx.TransactionType),
		Amount: &generated.Money{
			Amount:   tx.Amount.Amount.String(),
			Currency: tx.Amount.Currency,
		},
		Description:  tx.Description,
		MerchantName: &tx.MerchantName,
		PostingDate:  tx.PostingDate.Format(time.RFC3339),
		ValueDate:    tx.ValueDate.Format(time.RFC3339),
		AccountID:    tx.AccountID,
	}
	
	if tx.RunningBalance != nil {
		result.RunningBalance = &generated.Money{
			Amount:   tx.RunningBalance.Amount.String(),
			Currency: tx.RunningBalance.Currency,
		}
	}
	
	return result
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(time.RFC3339)
	return &formatted
}