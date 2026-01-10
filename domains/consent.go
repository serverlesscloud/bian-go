package domains

import (
	"context"

	"github.com/serverlesscloud/bian-go/models"
)

// ConsentService defines operations for consent management following BIAN Customer Consent Management service domain.
// This interface implements a subset of BIAN v13.0.0 operations focused on read-only consent retrieval.
//
// BIAN Alignment:
// - RetrieveConsent maps to BIAN "Retrieve Consent" operation
// - RetrieveConsentStatus maps to BIAN "Retrieve Consent Status" operation
type ConsentService interface {
	// RetrieveConsent retrieves full consent details by consent ID.
	// Returns the consent information or an error if the consent is not found.
	//
	// BIAN Operation: Retrieve Consent
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - consentID: Unique identifier for the consent
	//
	// Returns:
	//   - Consent details if found
	//   - Error if consent not found, access denied, or internal error
	RetrieveConsent(ctx context.Context, consentID string) (*models.Consent, error)

	// RetrieveConsentStatus retrieves only the status of a consent by consent ID.
	// Returns the consent status or an error if the consent is not found.
	//
	// BIAN Operation: Retrieve Consent Status
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - consentID: Unique identifier for the consent
	//
	// Returns:
	//   - Consent status if found
	//   - Error if consent not found, access denied, or internal error
	RetrieveConsentStatus(ctx context.Context, consentID string) (models.ConsentStatus, error)
}