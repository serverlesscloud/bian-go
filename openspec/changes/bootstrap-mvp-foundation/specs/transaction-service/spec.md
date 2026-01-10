# Transaction Service Capability

## ADDED Requirements

### Requirement: Transaction Retrieval
The system SHALL provide operations to retrieve payment transaction information following BIAN Payment Execution service domain (v13.0.0 subset). Transaction retrieval operations return canonical transaction models with BIAN-compliant data structures.

#### Scenario: Retrieve transaction by ID - Success
- **GIVEN** a transaction exists with ID "txn-98765"
- **WHEN** RetrievePaymentTransaction is called with transactionID "txn-98765"
- **THEN** the system returns Transaction model with ID, reference, type, amount, description, merchant, dates
- **AND** the transaction type is one of: Debit, Credit, Transfer, Payment, Fee
- **AND** the amount is represented as Money type with decimal precision

#### Scenario: Retrieve transaction by ID - Not Found
- **GIVEN** no transaction exists with ID "txn-00000"
- **WHEN** RetrievePaymentTransaction is called with transactionID "txn-00000"
- **THEN** the system returns an error indicating transaction not found
- **AND** the error message is user-friendly

### Requirement: Transaction History Query
The system SHALL provide operations to retrieve transaction history for a banking account with filtering and pagination support. Transaction history queries return ordered lists of Transaction models.

#### Scenario: Retrieve transaction history - Basic
- **GIVEN** an account "acc-12345" with 30 transactions
- **WHEN** RetrievePaymentTransactionHistory is called with accountID "acc-12345" and default options
- **THEN** the system returns a list of Transaction models ordered by posting date descending (newest first)
- **AND** each transaction includes amount, description, merchant, posting date, value date

#### Scenario: Retrieve transaction history - Date Range Filter
- **GIVEN** an account with transactions spanning 90 days
- **WHEN** RetrievePaymentTransactionHistory is called with date range filter (last 30 days)
- **THEN** the system returns only transactions within the specified date range
- **AND** transactions outside the range are excluded
- **AND** the list is ordered by posting date descending

#### Scenario: Retrieve transaction history - Pagination
- **GIVEN** an account with 30 transactions
- **WHEN** RetrievePaymentTransactionHistory is called with limit=10, offset=0
- **THEN** the system returns the first 10 transactions
- **WHEN** called again with limit=10, offset=10
- **THEN** the system returns the next 10 transactions
- **AND** no overlap with previous results

#### Scenario: Retrieve transaction history - Empty Result
- **GIVEN** an account with no transactions in the specified date range
- **WHEN** RetrievePaymentTransactionHistory is called with date range filter
- **THEN** the system returns an empty list (not an error)
- **AND** the response indicates success with zero results

### Requirement: Transaction Service Interface
The system SHALL define TransactionService as a Go interface with context-aware operations. The interface includes history options for filtering and pagination.

#### Scenario: Interface Contract
- **GIVEN** the TransactionService interface
- **WHEN** inspecting interface methods
- **THEN** RetrievePaymentTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) exists
- **AND** RetrievePaymentTransactionHistory(ctx context.Context, accountID string, opts HistoryOptions) ([]*models.Transaction, error) exists

#### Scenario: History Options Structure
- **GIVEN** the HistoryOptions struct
- **WHEN** inspecting struct fields
- **THEN** the struct includes FromDate (optional time.Time), ToDate (optional time.Time)
- **AND** the struct includes Limit (int), Offset (int) for pagination
- **AND** default values are sensible (Limit=100, Offset=0)

### Requirement: BIAN Alignment
The system SHALL align with BIAN Payment Execution service domain v13.0.0 specification for naming, semantics, and data structures. The MVP implements a subset of BIAN operations for read-only transaction queries.

#### Scenario: BIAN Operation Naming
- **GIVEN** BIAN Payment Execution service domain specification
- **WHEN** comparing operation names
- **THEN** RetrievePaymentTransaction matches BIAN naming pattern
- **AND** RetrievePaymentTransactionHistory aligns with BIAN query patterns
- **AND** documentation references BIAN v13.0.0

#### Scenario: Deferred Operations Documentation
- **GIVEN** BIAN Payment Execution has operations for payment initiation, execution, status tracking
- **WHEN** reviewing TransactionService interface
- **THEN** documentation clearly states implemented subset (read-only query operations)
- **AND** documentation lists deferred operations (InitiatePaymentExecution, UpdatePaymentExecution, RequestPaymentExecutionStatus, etc.)
- **AND** write operations are deferred to future changes

### Requirement: Transaction Model Canonical Structure
The system SHALL define Transaction model with BIAN-aligned fields representing payment transaction details. The model supports comprehensive transaction information including merchant details and running balance.

#### Scenario: Transaction Model Fields
- **GIVEN** the Transaction model
- **WHEN** inspecting model fields
- **THEN** the model includes ID (string), Reference (string), Type (TransactionType enum)
- **AND** the model includes Amount (Money), Description (string), MerchantName (optional string)
- **AND** the model includes PostingDate (time.Time), ValueDate (time.Time)
- **AND** the model includes RunningBalance (optional Money)

#### Scenario: Transaction Model JSON Serialization
- **GIVEN** a Transaction model with populated fields
- **WHEN** marshaling to JSON
- **THEN** JSON output includes all non-zero fields
- **AND** Money fields serialize as {"amount": "123.45", "currency": "AUD"}
- **AND** date fields use ISO 8601 format
- **AND** optional fields use omitempty tag

#### Scenario: Transaction Type Validation
- **GIVEN** a Transaction model
- **WHEN** setting transaction type
- **THEN** the type is validated against TransactionType enum (Debit, Credit, Transfer, Payment, Fee)
- **AND** invalid types are rejected

### Requirement: Running Balance Calculation
The system SHALL support optional running balance tracking in transaction history. Running balance represents account balance after each transaction when ordered chronologically.

#### Scenario: Running Balance Present
- **GIVEN** a transaction history with running balances
- **WHEN** querying transaction history
- **THEN** each transaction includes running balance as Money type
- **AND** running balances decrease for debits, increase for credits
- **AND** running balance currency matches account currency

#### Scenario: Running Balance Absent
- **GIVEN** a transaction without running balance information
- **WHEN** querying the transaction
- **THEN** the RunningBalance field is nil (not zero)
- **AND** JSON serialization omits the field
