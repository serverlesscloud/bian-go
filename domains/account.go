package domains

import (
	"context"

	"github.com/serverlesscloud/bian-go/models"
)

// AccountService defines operations for account management following BIAN Current Account Fulfillment service domain.
// This interface implements a subset of BIAN v13.0.0 operations focused on read-only account retrieval.
//
// BIAN Alignment:
// - RetrieveCurrentAccount maps to BIAN "Retrieve Current Account" operation
// - RetrieveCurrentAccountBalance maps to BIAN "Retrieve Current Account Balance" operation
type AccountService interface {
	// RetrieveCurrentAccount retrieves account details by account ID.
	// Returns the account information or an error if the account is not found or inaccessible.
	//
	// BIAN Operation: Retrieve Current Account
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - accountID: Unique identifier for the account
	//
	// Returns:
	//   - Account details if found
	//   - Error if account not found, access denied, or internal error
	RetrieveCurrentAccount(ctx context.Context, accountID string) (*models.Account, error)

	// RetrieveCurrentAccountBalance retrieves the current balance for an account.
	// Returns the current balance information or an error if the account is not found.
	//
	// BIAN Operation: Retrieve Current Account Balance
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - accountID: Unique identifier for the account
	//
	// Returns:
	//   - Balance information if account found
	//   - Error if account not found, access denied, or internal error
	RetrieveCurrentAccountBalance(ctx context.Context, accountID string) (*models.Balance, error)
}