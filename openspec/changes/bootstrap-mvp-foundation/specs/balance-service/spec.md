# Balance Service Capability

## ADDED Requirements

### Requirement: Account Balance Retrieval
The system SHALL provide operations to retrieve account balance information following BIAN Account Balance Management service domain (v13.0.0 subset). Balance retrieval operations return canonical balance models with multiple balance types (current, available, pending).

#### Scenario: Retrieve account balance - Single Balance
- **GIVEN** an account "acc-12345" with current balance $1,234.56 AUD
- **WHEN** RetrieveAccountBalance is called with accountID "acc-12345"
- **THEN** the system returns a list containing one Balance model
- **AND** the balance type is "Current"
- **AND** the amount is $1,234.56 AUD with decimal precision
- **AND** the balance includes timestamp indicating when it was calculated

#### Scenario: Retrieve account balance - Multiple Balance Types
- **GIVEN** an account with current balance $1,234.56 and available balance $1,100.00 AUD
- **WHEN** RetrieveAccountBalance is called
- **THEN** the system returns a list containing two Balance models
- **AND** one balance has type "Current" with amount $1,234.56
- **AND** one balance has type "Available" with amount $1,100.00
- **AND** available balance is less than or equal to current balance

#### Scenario: Retrieve account balance - Not Found
- **GIVEN** no account exists with ID "acc-99999"
- **WHEN** RetrieveAccountBalance is called with accountID "acc-99999"
- **THEN** the system returns an error indicating account not found
- **AND** the error message is user-friendly

#### Scenario: Retrieve account balance - Multi-Currency Account
- **GIVEN** an account with balances in multiple currencies (AUD and USD)
- **WHEN** RetrieveAccountBalance is called
- **THEN** the system returns Balance models for each currency
- **AND** each balance has correct currency code (ISO 4217)
- **AND** balances for different currencies are distinct entries

### Requirement: Balance Service Interface
The system SHALL define BalanceService as a Go interface with context-aware operations. Balance operations return lists of Balance models to support multiple balance types and currencies.

#### Scenario: Interface Contract
- **GIVEN** the BalanceService interface
- **WHEN** inspecting interface methods
- **THEN** RetrieveAccountBalance(ctx context.Context, accountID string) ([]*models.Balance, error) exists
- **AND** the method returns a slice of Balance pointers (not a single balance)
- **AND** the method accepts context for cancellation support

#### Scenario: Context Cancellation
- **GIVEN** a RetrieveAccountBalance operation is in progress
- **WHEN** the context is cancelled
- **THEN** the operation returns context.Canceled error
- **AND** no partial results are returned

### Requirement: BIAN Alignment
The system SHALL align with BIAN Account Balance Management service domain v13.0.0 specification for naming, semantics, and data structures. The MVP implements a subset of BIAN operations for read-only balance queries.

#### Scenario: BIAN Operation Naming
- **GIVEN** BIAN Account Balance Management service domain specification
- **WHEN** comparing operation names
- **THEN** RetrieveAccountBalance matches BIAN naming pattern
- **AND** operation semantics match BIAN specification
- **AND** documentation references BIAN v13.0.0

#### Scenario: Deferred Operations Documentation
- **GIVEN** BIAN Account Balance Management service domain
- **WHEN** reviewing BalanceService interface
- **THEN** documentation clearly states implemented subset (read-only query operations)
- **AND** documentation lists deferred operations (InitiateAccountBalanceCalculation, UpdateAccountBalance, etc.)
- **AND** write/update operations are deferred to future changes

### Requirement: Balance Model Canonical Structure
The system SHALL define Balance model with BIAN-aligned fields representing account balance information. The model supports multiple balance types and precise decimal amounts.

#### Scenario: Balance Model Fields
- **GIVEN** the Balance model
- **WHEN** inspecting model fields
- **THEN** the model includes Type (BalanceType enum)
- **AND** the model includes Amount (Money with decimal precision)
- **AND** the model includes Timestamp (time.Time)
- **AND** all fields are required (no optional fields)

#### Scenario: Balance Type Enum
- **GIVEN** the BalanceType enum
- **WHEN** inspecting enum values
- **THEN** the enum includes "Current" (ledger balance)
- **AND** the enum includes "Available" (current minus pending transactions)
- **AND** the enum includes "Pending" (sum of pending transactions)
- **AND** the enum supports JSON serialization (string format)

#### Scenario: Balance Model JSON Serialization
- **GIVEN** a Balance model with type "Current", amount $1,234.56 AUD, and timestamp
- **WHEN** marshaling to JSON
- **THEN** JSON output is {"type": "Current", "amount": {"amount": "1234.56", "currency": "AUD"}, "timestamp": "2026-01-10T10:30:00Z"}
- **AND** Money amount uses string format for precision
- **AND** timestamp uses ISO 8601 format

### Requirement: Balance Consistency Rules
The system SHALL enforce logical consistency rules for account balances. Available balance must not exceed current balance, and balance relationships must be maintained.

#### Scenario: Available Balance Constraint
- **GIVEN** an account with current balance and available balance
- **WHEN** comparing balance amounts
- **THEN** available balance is less than or equal to current balance
- **AND** the difference represents pending debits or holds

#### Scenario: Pending Balance Calculation
- **GIVEN** an account with current balance and available balance
- **WHEN** calculating pending transactions
- **THEN** pending amount equals current balance minus available balance
- **AND** pending amount is non-negative

#### Scenario: Zero Balance Handling
- **GIVEN** an account with zero balance
- **WHEN** RetrieveAccountBalance is called
- **THEN** the system returns Balance model with amount $0.00 (not nil)
- **AND** the balance type is correctly set
- **AND** currency code is present

### Requirement: Balance Timestamp Accuracy
The system SHALL include timestamps with all balance information to indicate when balances were calculated or retrieved. Timestamps support temporal queries and audit requirements.

#### Scenario: Balance Timestamp Present
- **GIVEN** a balance retrieval operation
- **WHEN** the operation completes successfully
- **THEN** the returned Balance model includes a timestamp
- **AND** the timestamp is in UTC
- **AND** the timestamp represents when the balance was calculated (not retrieved)

#### Scenario: Balance Timestamp Ordering
- **GIVEN** multiple balances for the same account
- **WHEN** comparing timestamps
- **THEN** all balances should have the same or very recent timestamps (within seconds)
- **AND** timestamps accurately reflect balance calculation time
