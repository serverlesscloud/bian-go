package domains

import (
	"context"

	"github.com/serverlesscloud/bian-go/models"
)

// BalanceService defines operations for balance management following BIAN Account Balance Management service domain.
// This interface implements a subset of BIAN v13.0.0 operations focused on read-only balance retrieval.
//
// BIAN Alignment:
// - RetrieveAccountBalance maps to BIAN "Retrieve Account Balance" operation
type BalanceService interface {
	// RetrieveAccountBalance retrieves all balance types (current, available, pending) for an account.
	// Returns a list of balance information or an error if the account is not found.
	//
	// BIAN Operation: Retrieve Account Balance
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - accountID: Unique identifier for the account
	//
	// Returns:
	//   - List of balance information (current, available, pending balances)
	//   - Error if account not found, access denied, or internal error
	RetrieveAccountBalance(ctx context.Context, accountID string) ([]*models.Balance, error)
}