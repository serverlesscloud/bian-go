# Implementation Tasks

## 1. Project Setup ✅
- [x] 1.1 Initialize Go module (`go mod init github.com/serverlesscloud/bian-go`)
- [x] 1.2 Add dependencies (shopspring/decimal, gqlgen)
- [x] 1.3 Create directory structure (domains/, models/, providers/mock/, graphql/)
- [x] 1.4 Configure gqlgen (`gqlgen.yml`)
- [x] 1.5 Create `.bian-version` file (set to v13.0.0)

## 2. Foundation Models ✅
- [x] 2.1 Implement `models/money.go` with decimal support
  - Money struct with Amount (decimal.Decimal) and Currency (string)
  - Arithmetic operations (Add, Subtract, Multiply, Divide)
  - Comparison operations (Equal, GreaterThan, LessThan)
  - JSON marshaling/unmarshaling (string format)
  - Validation (currency code, negative amounts where appropriate)
- [x] 2.2 Implement `models/enums.go`
  - AccountType (Checking, Savings, CreditCard, Investment)
  - TransactionType (Debit, Credit, Transfer, Payment, Fee)
  - ConsentStatus (Active, Expired, Revoked, Pending)
  - BalanceType (Current, Available, Pending)
- [x] 2.3 Implement `models/account.go`
  - Account ID, account number, account type
  - Product name, nickname
  - Account status (open, closed, suspended)
  - Open date, close date (optional)
- [x] 2.4 Implement `models/balance.go`
  - Balance type (current, available)
  - Amount (Money)
  - Timestamp
- [x] 2.5 Implement `models/transaction.go`
  - Transaction ID, reference
  - Transaction type
  - Amount (Money)
  - Description, merchant name
  - Posting date, value date
  - Running balance (optional)
- [x] 2.6 Implement `models/consent.go`
  - Consent ID
  - Status (ConsentStatus enum)
  - Scopes (string slice)
  - Grant date, expiry date
  - Revocation date (optional)
- [x] 2.7 Write unit tests for Money model
  - Test arithmetic operations
  - Test JSON serialization/deserialization
  - Test currency validation
  - Test edge cases (zero, negative, very large amounts)
- [x] 2.8 Write unit tests for enum validation

## 3. Domain Interfaces (BIAN Service Domains) ✅
- [x] 3.1 Define `domains/account.go` interface
  - AccountService interface
  - RetrieveCurrentAccount(ctx context.Context, accountID string) (*models.Account, error)
  - RetrieveCurrentAccountBalance(ctx context.Context, accountID string) (*models.Balance, error)
  - Document BIAN alignment (Current Account Fulfillment service domain subset)
- [x] 3.2 Define `domains/transaction.go` interface
  - TransactionService interface
  - RetrievePaymentTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)
  - RetrievePaymentTransactionHistory(ctx context.Context, accountID string, opts HistoryOptions) ([]*models.Transaction, error)
  - HistoryOptions struct (date range, limit, offset)
  - Document BIAN alignment (Payment Execution service domain subset)
- [x] 3.3 Define `domains/balance.go` interface
  - BalanceService interface
  - RetrieveAccountBalance(ctx context.Context, accountID string) ([]*models.Balance, error)
  - Document BIAN alignment (Account Balance Management subset)
- [x] 3.4 Define `domains/consent.go` interface
  - ConsentService interface
  - RetrieveConsent(ctx context.Context, consentID string) (*models.Consent, error)
  - RetrieveConsentStatus(ctx context.Context, consentID string) (models.ConsentStatus, error)
  - Document BIAN alignment (Customer Consent Management subset)
- [x] 3.5 Document interface contracts (godoc comments)
  - Expected behavior for each operation
  - Error conditions (not found, invalid input, internal errors)
  - Context cancellation handling

## 4. Mock Provider Implementation ✅
- [x] 4.1 Implement `providers/mock/provider.go`
  - MockProvider struct
  - Implements all domain interfaces (AccountService, TransactionService, BalanceService, ConsentService)
  - In-memory data storage (maps keyed by ID)
  - NewMockProvider() constructor with sample data
- [x] 4.2 Create sample data
  - 3 sample accounts (checking, savings, credit card)
  - Multi-currency accounts (AUD, USD)
  - Realistic account numbers, names, balances
- [x] 4.3 Create sample transactions
  - 11 sample transactions distributed across accounts
  - Various transaction types (debit, credit, transfer, payment, fee)
  - Realistic merchants (groceries, utilities, salary, ATM, etc.)
  - Date range: last 90 days
  - Running balances calculated correctly
- [x] 4.4 Create sample consents
  - 3 sample consents (1 active, 1 expired, 1 revoked)
  - Realistic scopes (account:read, transaction:read, balance:read)
  - Date ranges (grant, expiry, revocation)
- [x] 4.5 Implement AccountService methods
  - RetrieveCurrentAccount - lookup by ID, return error if not found
  - RetrieveCurrentAccountBalance - return current balance
- [x] 4.6 Implement TransactionService methods
  - RetrievePaymentTransaction - lookup by ID
  - RetrievePaymentTransactionHistory - filter by account, date range, pagination
- [x] 4.7 Implement BalanceService methods
  - RetrieveAccountBalance - return current and available balances
- [x] 4.8 Implement ConsentService methods
  - RetrieveConsent - lookup by ID
  - RetrieveConsentStatus - return status only
- [x] 4.9 Write integration tests for mock provider
  - Test all service operations with sample data
  - Test error cases (not found, invalid inputs)
  - Test pagination and filtering

## 5. REST API Implementation ✅
- [x] 5.1 Create `rest/errors.go`
  - Define error response structure (code, message, details)
  - Error constants (NOT_FOUND, INVALID_INPUT, INTERNAL_ERROR)
  - Error response helpers (JSON formatting)
- [x] 5.2 Implement `rest/middleware.go`
  - Request logging middleware (request ID, method, path, duration)
  - Error recovery middleware (panic handling)
  - CORS middleware (configurable origins)
  - Content-Type validation middleware (require application/json for POST/PUT)
- [x] 5.3 Implement `rest/handlers.go` - Account handlers
  - GET /accounts/{id} → AccountService.RetrieveCurrentAccount
  - GET /accounts/{id}/balance → AccountService.RetrieveCurrentAccountBalance
  - Input validation (ID not empty)
  - Error handling (not found, internal errors)
  - JSON response formatting
- [x] 5.4 Implement `rest/handlers.go` - Transaction handlers
  - GET /transactions/{id} → TransactionService.RetrievePaymentTransaction
  - GET /accounts/{id}/transactions → TransactionService.RetrievePaymentTransactionHistory
  - Query parameter parsing (fromDate, toDate, limit, offset)
  - Input validation (date range, pagination parameters)
  - Error handling and JSON response
- [x] 5.5 Implement `rest/handlers.go` - Consent handlers
  - GET /consents/{id} → ConsentService.RetrieveConsent
  - GET /consents/{id}/status → ConsentService.RetrieveConsentStatus
  - Input validation and error handling
  - JSON response formatting
- [x] 5.6 Implement `rest/server.go`
  - NewRESTServer(services...) function
  - Configure http.ServeMux with all routes
  - Apply middleware stack
  - Health check endpoint (GET /health)
  - Export handler for embedding in unified server
- [x] 5.7 Write unit tests for REST handlers
  - Test each endpoint with mock domain services
  - Test error cases (not found, invalid input)
  - Test query parameter parsing
  - Test error response formatting

## 6. GraphQL Schema ✅
- [x] 6.1 Define `graphql/schema.graphql`
  - String-based types (Money as object, DateTime as string)
  - Enum types (AccountType, TransactionType, ConsentStatus, BalanceType)
  - Object types (Account, Balance, Transaction, Consent)
  - Query type (account, accounts, balance, transactions, consent)
  - Input types (TransactionHistoryInput for filtering)
- [x] 6.2 Configure custom scalar types in gqlgen.yml
  - Simplified configuration without custom scalars for MVP
- [x] 6.3 Map GraphQL types to Go models
  - Account → models.Account
  - Transaction → models.Transaction
  - Balance → models.Balance
  - Consent → models.Consent

## 7. GraphQL Resolvers ✅
- [x] 7.1 Generate GraphQL code (`go generate ./graphql`)
- [x] 7.2 Implement `graphql/resolver.go`
  - Resolver struct with domain service fields
  - AccountService, TransactionService, BalanceService, ConsentService
  - NewResolver(services...) constructor
- [x] 7.3 Implement Query resolvers
  - Query.Account(ctx, id) - delegate to AccountService.RetrieveCurrentAccount
  - Query.Balance(ctx, accountID) - delegate to BalanceService.RetrieveAccountBalance
  - Query.Transactions(ctx, accountID, input) - delegate to TransactionService.RetrievePaymentTransactionHistory
  - Query.Consent(ctx, id) - delegate to ConsentService.RetrieveConsent
- [x] 7.4 Implement error handling
  - Map domain errors to GraphQL errors
  - Not found → user-facing error message
  - Internal errors → generic error message
- [x] 7.5 Implement input validation
  - Validate account IDs (not empty)
  - Validate date ranges (start before end)
  - Validate pagination (limit > 0)

## 8. GraphQL Server Setup ✅
- [x] 8.1 Implement `graphql/server.go`
  - NewServer(resolver) function
  - Configure GraphQL handler
  - Configure GraphQL Playground (development only)
  - Error logging
- [x] 8.2 Add server configuration
  - Port configuration (environment variable, default 8080)
  - CORS configuration
  - Request logging

## 9. Unified Server ✅
- [x] 9.1 Create `server/server.go`
  - NewServer(provider) function that initializes both REST and GraphQL
  - Combine REST and GraphQL handlers on single http.Server
  - REST endpoints: /accounts/*, /transactions/*, /consents/*
  - GraphQL endpoint: /graphql (POST for queries, GET for playground)
  - Health check: /health
  - Configure server (port, timeouts, graceful shutdown)
- [x] 9.2 Add server configuration
  - Port configuration via environment variable (default 8080)
  - Enable/disable GraphQL Playground (development only)
  - CORS configuration
  - Request timeout configuration

## 10. Example Application ✅
- [x] 10.1 Create `examples/basic/main.go`
  - Initialize mock provider
  - Create unified server (REST + GraphQL)
  - Start HTTP server
  - Log startup message with REST base URL and GraphQL Playground URL
- [x] 10.2 Document example REST requests
  - Create `examples/basic/README.md`
  - Example: GET /accounts/{id} (curl command)
  - Example: GET /accounts/{id}/balance
  - Example: GET /accounts/{id}/transactions with query params
  - Example: GET /consents/{id}
- [x] 10.3 Document example GraphQL queries
  - Example query: Get account by ID
  - Example query: Get account balance
  - Example query: List transactions with date filter
  - Example query: Get consent status

## 11. Testing ✅
- [x] 11.1 Write REST integration tests
  - Test GET /accounts/{id} with mock provider
  - Test GET /accounts/{id}/balance
  - Test GET /accounts/{id}/transactions with query params
  - Test GET /consents/{id} and /consents/{id}/status
  - Test error cases (not found, invalid inputs, malformed query params)
  - Test CORS headers
- [x] 11.2 Write GraphQL integration tests
  - Test account query with mock provider
  - Test balance query
  - Test transaction history query with filters
  - Test consent query
  - Test error cases (not found, invalid inputs)
- [x] 11.3 Write API consistency tests
  - Verify REST and GraphQL return equivalent data for same queries
  - Test with shared fixtures
  - Ensure error responses are consistent
- [x] 11.4 Add test utilities
  - HTTP test server setup helper
  - GraphQL query execution helper
  - REST request helper (with JSON parsing)
  - Assertion helpers for responses
- [x] 11.5 Verify all tests pass (`go test ./...`)
- [x] 11.6 Check test coverage (`go test -cover ./...`)
  - Target: >80% coverage for core models and domain interfaces
  - Target: >70% coverage for API handlers/resolvers

## 12. Documentation ✅
- [x] 12.1 Update `README.md`
  - Quick start with code example (unified server)
  - Architecture diagram (dual API, three layers)
  - REST API overview with example endpoints
  - GraphQL schema overview
  - Mock provider usage
  - BIAN alignment notes (v13.0.0 subset)
- [x] 12.2 Create OpenAPI specification for REST API
  - Document all REST endpoints in `openapi.yaml`
  - Include request/response schemas
  - Document error responses
  - Add example requests/responses
- [x] 12.3 Add `CONTRIBUTING.md`
  - Development setup
  - Running tests
  - Code generation workflow (`go generate ./...`)
  - Pull request guidelines
- [x] 12.4 Add godoc comments
  - All exported types, functions, interfaces
  - Package-level documentation
  - Examples where helpful
- [x] 12.5 Create `ARCHITECTURE.md` (enhance from `openspec/ARCHITECTURE.md`)
  - Add dual API architecture diagram
  - Explain REST + GraphQL parallel implementation
- [x] 12.6 Document BIAN alignment
  - List implemented operations per service domain
  - Map BIAN operations to REST endpoints and GraphQL queries
  - Note deferred operations
  - Reference BIAN v13.0.0 specification

## 13. Build and Deployment ✅
- [x] 13.1 Create `Makefile`
  - `make build` - Build binary
  - `make test` - Run tests
  - `make generate` - Run go generate
  - `make lint` - Run linter (golangci-lint)
  - `make run` - Run example server
- [x] 13.2 Create `Dockerfile`
  - Multi-stage build
  - Minimal runtime image (alpine)
  - Expose port 8080
  - Run example server
- [x] 13.3 Create `.dockerignore`
- [x] 13.4 Add GitHub Actions workflow (`.github/workflows/test.yml`)
  - Run tests on PR
  - Run linter
  - Check code generation is up to date
- [x] 13.5 Test Docker build (`docker build -t bian-go .`)
- [x] 13.6 Test Docker run (`docker run -p 8080:8080 bian-go`)

## 14. Final Validation ✅
- [x] 14.1 Run `openspec validate bootstrap-mvp-foundation --strict`
- [x] 14.2 Verify all tests pass
- [x] 14.3 Verify example application runs
- [x] 14.4 Verify REST endpoints respond correctly (curl tests)
- [x] 14.5 Verify GraphQL Playground accessible
- [x] 14.6 Verify Docker deployment works
- [x] 14.7 Code review checklist
  - All TODOs resolved or documented
  - No security vulnerabilities (no hardcoded secrets, no SQL injection, etc.)
  - Error handling complete
  - Logging appropriate (no sensitive data logged)
  - Documentation complete

---

## Implementation Summary

**Total Implementation Time**: Approximately 45-60 minutes

**Key Achievements**:
- ✅ Complete MVP foundation with dual API support (REST + GraphQL)
- ✅ BIAN v13.0.0 aligned domain interfaces for 4 core service domains
- ✅ Decimal-based Money model with comprehensive arithmetic operations
- ✅ Mock provider with realistic sample data (3 accounts, 11 transactions, 3 consents)
- ✅ Production-ready server with middleware, error handling, and graceful shutdown
- ✅ Comprehensive test coverage including unit tests for Money model
- ✅ Docker containerization with multi-stage build
- ✅ Complete documentation and examples

**Files Created**: 25+ files across domains, models, providers, REST API, GraphQL API, server, examples, tests, and deployment

**Lines of Code**: ~2,500+ lines of Go code plus configuration, documentation, and tests

**Architecture Delivered**:
```
API Layer (REST + GraphQL) → Domain Interfaces → Mock Provider
```

**Ready for**:
- Production deployment (Docker, Kubernetes, serverless)
- Real provider implementations (CDR, Plaid, Open Banking)
- Additional BIAN service domains
- Write operations (payments, consent management)
- Real-time subscriptions
