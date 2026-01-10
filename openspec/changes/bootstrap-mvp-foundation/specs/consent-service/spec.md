# Consent Service Capability

## ADDED Requirements

### Requirement: Consent Retrieval
The system SHALL provide operations to retrieve OAuth consent information following BIAN Customer Consent Management service domain (v13.0.0 subset). Consent retrieval operations return canonical consent models with BIAN-compliant data structures.

#### Scenario: Retrieve consent by ID - Success
- **GIVEN** a consent exists with ID "consent-abc123"
- **WHEN** RetrieveConsent is called with consentID "consent-abc123"
- **THEN** the system returns Consent model with ID, status, scopes, grant date, expiry date
- **AND** the consent status is one of: Active, Expired, Revoked, Pending
- **AND** scopes are returned as a list of strings (e.g., ["account:read", "transaction:read"])

#### Scenario: Retrieve consent by ID - Not Found
- **GIVEN** no consent exists with ID "consent-xyz999"
- **WHEN** RetrieveConsent is called with consentID "consent-xyz999"
- **THEN** the system returns an error indicating consent not found
- **AND** the error message is user-friendly

#### Scenario: Retrieve consent by ID - Invalid Input
- **GIVEN** an empty consent ID
- **WHEN** RetrieveConsent is called with consentID ""
- **THEN** the system returns a validation error
- **AND** the error message indicates invalid consent ID

### Requirement: Consent Status Query
The system SHALL provide operations to retrieve consent status without full consent details. Status queries support quick consent validation checks.

#### Scenario: Retrieve consent status - Active
- **GIVEN** a consent with status "Active" and expiry date in the future
- **WHEN** RetrieveConsentStatus is called
- **THEN** the system returns ConsentStatus enum value "Active"
- **AND** no other consent details are returned (status only)

#### Scenario: Retrieve consent status - Expired
- **GIVEN** a consent with status "Expired" and expiry date in the past
- **WHEN** RetrieveConsentStatus is called
- **THEN** the system returns ConsentStatus enum value "Expired"

#### Scenario: Retrieve consent status - Revoked
- **GIVEN** a consent with status "Revoked" and revocation date present
- **WHEN** RetrieveConsentStatus is called
- **THEN** the system returns ConsentStatus enum value "Revoked"

#### Scenario: Retrieve consent status - Not Found
- **GIVEN** no consent exists with the given ID
- **WHEN** RetrieveConsentStatus is called
- **THEN** the system returns an error indicating consent not found

### Requirement: Consent Service Interface
The system SHALL define ConsentService as a Go interface with context-aware operations. All operations accept context.Context for cancellation and timeout support.

#### Scenario: Interface Contract
- **GIVEN** the ConsentService interface
- **WHEN** inspecting interface methods
- **THEN** RetrieveConsent(ctx context.Context, consentID string) (*models.Consent, error) exists
- **AND** RetrieveConsentStatus(ctx context.Context, consentID string) (models.ConsentStatus, error) exists

#### Scenario: Context Cancellation
- **GIVEN** a RetrieveConsent operation is in progress
- **WHEN** the context is cancelled
- **THEN** the operation returns context.Canceled error
- **AND** no partial results are returned

### Requirement: BIAN Alignment
The system SHALL align with BIAN Customer Consent Management service domain v13.0.0 specification for naming, semantics, and data structures. The MVP implements a subset of BIAN operations for read-only consent queries.

#### Scenario: BIAN Operation Naming
- **GIVEN** BIAN Customer Consent Management service domain specification
- **WHEN** comparing operation names
- **THEN** RetrieveConsent matches BIAN naming pattern
- **AND** operation semantics match BIAN specification
- **AND** documentation references BIAN v13.0.0

#### Scenario: Deferred Operations Documentation
- **GIVEN** BIAN Customer Consent Management has operations for consent creation, update, revocation
- **WHEN** reviewing ConsentService interface
- **THEN** documentation clearly states implemented subset (read-only query operations)
- **AND** documentation lists deferred operations (InitiateCustomerConsent, UpdateCustomerConsent, RevokeCustomerConsent, etc.)
- **AND** write operations are deferred to future changes

### Requirement: Consent Model Canonical Structure
The system SHALL define Consent model with BIAN-aligned fields representing OAuth consent details. The model supports comprehensive consent lifecycle tracking including dates and scope information.

#### Scenario: Consent Model Fields
- **GIVEN** the Consent model
- **WHEN** inspecting model fields
- **THEN** the model includes ID (string), Status (ConsentStatus enum)
- **AND** the model includes Scopes ([]string - list of permission scopes)
- **AND** the model includes GrantDate (time.Time), ExpiryDate (time.Time)
- **AND** the model includes RevocationDate (optional time.Time)

#### Scenario: Consent Model JSON Serialization
- **GIVEN** a Consent model with populated fields
- **WHEN** marshaling to JSON
- **THEN** JSON output includes id, status, scopes, grantDate, expiryDate
- **AND** date fields use ISO 8601 format
- **AND** scopes are JSON array of strings
- **AND** revocationDate is omitted if null (omitempty tag)

#### Scenario: Consent Status Enum
- **GIVEN** the ConsentStatus enum
- **WHEN** inspecting enum values
- **THEN** the enum includes "Active" (consent is valid and not expired)
- **AND** the enum includes "Expired" (consent has passed expiry date)
- **AND** the enum includes "Revoked" (consent was explicitly revoked)
- **AND** the enum includes "Pending" (consent awaiting user authorization)

### Requirement: Consent Scope Format
The system SHALL represent consent scopes as a list of strings following OAuth 2.0 scope conventions. Scopes define what data and operations the consent permits.

#### Scenario: Scope String Format
- **GIVEN** a Consent model with scopes
- **WHEN** inspecting scope values
- **THEN** scopes follow pattern "resource:permission" (e.g., "account:read")
- **AND** common scopes include "account:read", "transaction:read", "balance:read"
- **AND** scopes are case-sensitive
- **AND** invalid scope formats are validated

#### Scenario: Multiple Scopes
- **GIVEN** a consent with multiple permissions
- **WHEN** retrieving consent
- **THEN** the scopes list includes all granted scopes
- **AND** scopes are in consistent order (sorted alphabetically recommended)
- **AND** duplicate scopes are not present

#### Scenario: Empty Scopes
- **GIVEN** a consent with no scopes (invalid state)
- **WHEN** validating consent
- **THEN** the system treats this as invalid consent
- **AND** at least one scope is required for valid consent

### Requirement: Consent Lifecycle Dates
The system SHALL track consent lifecycle through grant date, expiry date, and optional revocation date. Date comparisons determine consent validity.

#### Scenario: Grant Date
- **GIVEN** a Consent model
- **WHEN** inspecting grant date
- **THEN** grant date represents when user authorized the consent
- **AND** grant date is always present (required field)
- **AND** grant date is in the past or present (not future)

#### Scenario: Expiry Date
- **GIVEN** a Consent model
- **WHEN** inspecting expiry date
- **THEN** expiry date represents when consent becomes invalid
- **AND** expiry date is always present (required field)
- **AND** expiry date is after grant date

#### Scenario: Revocation Date
- **GIVEN** a revoked consent
- **WHEN** inspecting revocation date
- **THEN** revocation date is present and represents when user revoked consent
- **AND** revocation date is between grant date and expiry date (typically)
- **AND** revocation date is in the past or present (not future)

#### Scenario: Revocation Date Absent
- **GIVEN** an active or expired consent that was never revoked
- **WHEN** inspecting revocation date
- **THEN** revocation date is nil (not zero time)
- **AND** JSON serialization omits the field

### Requirement: Consent Status Determination
The system SHALL determine consent status based on current time, expiry date, and revocation date. Status determination follows OAuth consent lifecycle rules.

#### Scenario: Active Consent
- **GIVEN** current time is after grant date
- **AND** current time is before expiry date
- **AND** consent has not been revoked
- **WHEN** checking consent status
- **THEN** the status is "Active"

#### Scenario: Expired Consent
- **GIVEN** current time is after expiry date
- **AND** consent has not been revoked
- **WHEN** checking consent status
- **THEN** the status is "Expired"
- **AND** the consent cannot be used for authorization

#### Scenario: Revoked Consent
- **GIVEN** consent has revocation date in the past
- **WHEN** checking consent status
- **THEN** the status is "Revoked"
- **AND** the consent cannot be used for authorization
- **AND** revoked status takes precedence over expired status

#### Scenario: Pending Consent
- **GIVEN** consent exists but grant date is in the future (edge case)
- **OR** consent is awaiting user action
- **WHEN** checking consent status
- **THEN** the status is "Pending"
- **AND** the consent cannot be used for authorization yet

### Requirement: Consent Validation
The system SHALL validate consent data integrity including date relationships, scope format, and status consistency. Invalid consents are rejected at creation or update.

#### Scenario: Date Relationship Validation
- **GIVEN** a consent being validated
- **WHEN** checking date relationships
- **THEN** grant date must be before or equal to expiry date
- **AND** if revocation date exists, it must be after grant date
- **AND** if revocation date exists, status must be "Revoked"

#### Scenario: Status Consistency Validation
- **GIVEN** a consent with status "Revoked"
- **WHEN** validating consent
- **THEN** revocation date must be present
- **GIVEN** a consent with status "Active" or "Expired"
- **THEN** revocation date must be absent

#### Scenario: Scope Validation
- **GIVEN** a consent with scopes
- **WHEN** validating scopes
- **THEN** at least one scope must be present
- **AND** all scopes must follow "resource:permission" format
- **AND** no duplicate scopes are allowed
