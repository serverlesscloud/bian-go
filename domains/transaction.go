package domains

import (
	"context"
	"time"

	"github.com/serverlesscloud/bian-go/models"
)

// HistoryOptions provides filtering and pagination options for transaction history queries
type HistoryOptions struct {
	// Date range filtering
	FromDate *time.Time `json:"fromDate,omitempty"`
	ToDate   *time.Time `json:"toDate,omitempty"`
	
	// Pagination
	Limit  int `json:"limit,omitempty"`  // Maximum number of transactions to return (default: 50, max: 500)
	Offset int `json:"offset,omitempty"` // Number of transactions to skip (for pagination)
}

// TransactionService defines operations for transaction management following BIAN Payment Execution service domain.
// This interface implements a subset of BIAN v13.0.0 operations focused on read-only transaction retrieval.
//
// BIAN Alignment:
// - RetrievePaymentTransaction maps to BIAN "Retrieve Payment Transaction" operation
// - RetrievePaymentTransactionHistory maps to BIAN "Retrieve Payment Transaction History" operation
type TransactionService interface {
	// RetrievePaymentTransaction retrieves a specific transaction by transaction ID.
	// Returns the transaction details or an error if the transaction is not found.
	//
	// BIAN Operation: Retrieve Payment Transaction
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - transactionID: Unique identifier for the transaction
	//
	// Returns:
	//   - Transaction details if found
	//   - Error if transaction not found, access denied, or internal error
	RetrievePaymentTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)

	// RetrievePaymentTransactionHistory retrieves transaction history for an account with optional filtering.
	// Returns a list of transactions matching the criteria or an error if the account is not found.
	//
	// BIAN Operation: Retrieve Payment Transaction History
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - accountID: Unique identifier for the account
	//   - opts: Filtering and pagination options (date range, limit, offset)
	//
	// Returns:
	//   - List of transactions matching criteria (may be empty)
	//   - Error if account not found, invalid parameters, or internal error
	RetrievePaymentTransactionHistory(ctx context.Context, accountID string, opts HistoryOptions) ([]*models.Transaction, error)
}