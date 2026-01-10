# Account Service Capability

## ADDED Requirements

### Requirement: Account Retrieval
The system SHALL provide operations to retrieve banking account information following BIAN Current Account Fulfillment service domain (v13.0.0 subset). Account retrieval operations return canonical account models with BIAN-compliant data structures.

#### Scenario: Retrieve account by ID - Success
- **GIVEN** an account exists with ID "acc-12345"
- **WHEN** RetrieveCurrentAccount is called with accountID "acc-12345"
- **THEN** the system returns Account model with ID, account number, type, product name, status, and dates
- **AND** the account type is one of: Checking, Savings, CreditCard, Investment
- **AND** the account status is one of: Open, Closed, Suspended

#### Scenario: Retrieve account by ID - Not Found
- **GIVEN** no account exists with ID "acc-99999"
- **WHEN** RetrieveCurrentAccount is called with accountID "acc-99999"
- **THEN** the system returns an error indicating account not found
- **AND** the error message is user-friendly

#### Scenario: Retrieve account by ID - Invalid Input
- **GIVEN** an empty account ID
- **WHEN** RetrieveCurrentAccount is called with accountID ""
- **THEN** the system returns a validation error
- **AND** the error message indicates invalid account ID

### Requirement: Account Balance Query
The system SHALL provide operations to retrieve current and available balance for a banking account. Balance queries return Balance models with Money amounts using decimal precision.

#### Scenario: Retrieve account balance - Success
- **GIVEN** an account exists with ID "acc-12345" and current balance $1,234.56 AUD
- **WHEN** RetrieveCurrentAccountBalance is called with accountID "acc-12345"
- **THEN** the system returns a Balance model with amount $1,234.56 AUD
- **AND** the balance type is "Current"
- **AND** the amount uses decimal.Decimal type (not float64)
- **AND** the currency code is ISO 4217 compliant

#### Scenario: Retrieve account balance - Multi-Currency
- **GIVEN** an account with balances in multiple currencies (AUD, USD)
- **WHEN** RetrieveCurrentAccountBalance is called
- **THEN** the system returns balances for all currencies
- **AND** each balance has correct currency code

### Requirement: Account Service Interface
The system SHALL define AccountService as a Go interface with context-aware operations. All operations accept context.Context as the first parameter for cancellation and timeout support.

#### Scenario: Context Cancellation
- **GIVEN** a RetrieveCurrentAccount operation is in progress
- **WHEN** the context is cancelled
- **THEN** the operation returns context.Canceled error
- **AND** no partial results are returned

#### Scenario: Interface Contract
- **GIVEN** the AccountService interface
- **WHEN** inspecting interface methods
- **THEN** RetrieveCurrentAccount(ctx context.Context, accountID string) (*models.Account, error) exists
- **AND** RetrieveCurrentAccountBalance(ctx context.Context, accountID string) (*models.Balance, error) exists

### Requirement: BIAN Alignment
The system SHALL align with BIAN Current Account Fulfillment service domain v13.0.0 specification for naming, semantics, and data structures. The MVP implements a subset of BIAN operations, not the complete service domain.

#### Scenario: BIAN Operation Naming
- **GIVEN** BIAN Current Account service domain specification
- **WHEN** comparing operation names
- **THEN** RetrieveCurrentAccount matches BIAN naming pattern
- **AND** operation semantics match BIAN specification
- **AND** documentation references BIAN v13.0.0

#### Scenario: Deferred Operations Documentation
- **GIVEN** BIAN Current Account has 30+ operations
- **WHEN** reviewing AccountService interface
- **THEN** documentation clearly states implemented subset (2 operations)
- **AND** documentation lists deferred operations (InitiateCurrentAccountFacility, UpdateCurrentAccountFacility, RegisterDirectDebit, etc.)
- **AND** upgrade path is documented

### Requirement: Account Model Canonical Structure
The system SHALL define Account model with BIAN-aligned fields representing banking account details. The model supports JSON serialization for API responses.

#### Scenario: Account Model Fields
- **GIVEN** the Account model
- **WHEN** inspecting model fields
- **THEN** the model includes ID (string), AccountNumber (string), Type (AccountType enum)
- **AND** the model includes ProductName (string), Nickname (optional string)
- **AND** the model includes Status (string), OpenDate (time.Time), CloseDate (optional time.Time)

#### Scenario: Account Model JSON Serialization
- **GIVEN** an Account model with populated fields
- **WHEN** marshaling to JSON
- **THEN** JSON output includes all non-zero fields
- **AND** date fields use ISO 8601 format
- **AND** optional fields are omitted if empty (omitempty tag)

#### Scenario: Account Type Validation
- **GIVEN** an Account model
- **WHEN** setting account type
- **THEN** the type is validated against AccountType enum (Checking, Savings, CreditCard, Investment)
- **AND** invalid types are rejected
