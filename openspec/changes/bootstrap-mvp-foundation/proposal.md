# Change: Bootstrap MVP Foundation with Core BIAN Service Domains

## Why

This change establishes the foundational architecture for bian-go by implementing essential BIAN service domains required for any banking application. The MVP foundation provides the minimum viable set of interfaces and models needed to build fintech applications with open banking integration, starting with read-only operations and a mock provider for rapid development and testing.

Without this foundation, the project cannot demonstrate core value proposition: BIAN-compliant banking interfaces with dual API layers (REST and GraphQL) and provider-agnostic implementation.

## What Changes

This change introduces five core capabilities following a subset of BIAN v13 Service Landscape specification:

### New Service Domains (BIAN-aligned)
- **Account Service** - Current Account Fulfillment service domain (read-only subset)
  - `RetrieveCurrentAccount` - Get account details
  - `RetrieveCurrentAccountBalance` - Query balance information

- **Transaction Service** - Payment Execution service domain (read-only subset)
  - `RetrievePaymentTransaction` - Get transaction details
  - `RetrievePaymentTransactionHistory` - List transaction history

- **Balance Service** - Account Balance Management (read-only subset)
  - `RetrieveAccountBalance` - Query current and available balance
  - Support for multi-currency balances

- **Consent Service** - Customer Consent Management (read-only subset)
  - `RetrieveConsent` - Get consent details
  - `RetrieveConsentStatus` - Query consent state and expiry

### Foundation Models
- **Money Model** - Decimal-based currency representation
  - Uses `decimal.Decimal` type to avoid floating-point errors
  - ISO 4217 currency codes (AUD, USD, GBP, EUR, CAD)
  - Arithmetic operations (add, subtract, multiply, divide)
  - JSON serialization (string format for precision)

### Architecture
- **Three-layer design**: API Layer (REST + GraphQL) → Domain Interfaces → Providers
- **Dual API support**: Both REST and GraphQL APIs call domain interfaces directly (parallel implementations)
- **Provider abstraction**: All service domains implemented as Go interfaces
- **Mock provider**: In-memory implementation with realistic test data
- **REST API**: Standard RESTful endpoints using Go 1.22+ http.ServeMux (zero dependencies)
- **GraphQL API**: Query-only operations using gqlgen (no mutations or subscriptions in MVP)

### Testing Strategy
- Unit tests for domain interfaces (mock implementations)
- Unit tests for Money model (arithmetic, serialization, validation)
- Integration tests for REST endpoints against mock provider
- Integration tests for GraphQL queries against mock provider
- No external API dependencies required for MVP

## Impact

### Affected Specs
- **NEW**: `specs/account-service/` - Account retrieval capability
- **NEW**: `specs/transaction-service/` - Transaction query capability
- **NEW**: `specs/balance-service/` - Balance information capability
- **NEW**: `specs/money-model/` - Financial amount representation
- **NEW**: `specs/consent-service/` - OAuth consent management capability

### Affected Code
- **NEW**: `domains/account.go` - AccountService interface
- **NEW**: `domains/transaction.go` - TransactionService interface
- **NEW**: `domains/balance.go` - BalanceService interface
- **NEW**: `domains/consent.go` - ConsentService interface
- **NEW**: `models/money.go` - Money type with decimal support
- **NEW**: `models/account.go` - Account canonical model
- **NEW**: `models/transaction.go` - Transaction canonical model
- **NEW**: `models/balance.go` - Balance canonical model
- **NEW**: `models/consent.go` - Consent canonical model
- **NEW**: `models/enums.go` - Account types, transaction types, consent status
- **NEW**: `providers/mock/provider.go` - Mock implementation
- **NEW**: `rest/handlers.go` - REST HTTP handlers
- **NEW**: `rest/middleware.go` - REST middleware (logging, error handling)
- **NEW**: `rest/server.go` - REST server setup with http.ServeMux
- **NEW**: `graphql/schema.graphql` - GraphQL schema
- **NEW**: `graphql/resolver.go` - GraphQL resolvers
- **NEW**: `graphql/server.go` - GraphQL server setup

### Dependencies
- `github.com/shopspring/decimal` - Decimal arithmetic for Money model
- `github.com/99designs/gqlgen` - GraphQL server and code generation

### Migration Impact
None - this is the initial implementation with no existing code to migrate.

### Deployment Impact
- Containerized deployment (Docker, Kubernetes, serverless)
- No external dependencies required for mock provider
- REST API exposed on configurable port (default: 8080)
- GraphQL API exposed on same port at /graphql endpoint
- Both APIs can run on same server or separately

### Future Extensions
This MVP foundation enables future incremental additions:
- Write operations (payment initiation, consent creation)
- Real-time subscriptions (balance updates, transaction notifications)
- Additional providers (CDR, Plaid, Open Banking UK, PSD2)
- Additional service domains (Customer Information, Card Management)
- Full BIAN operation coverage (from subset to complete service domains)
