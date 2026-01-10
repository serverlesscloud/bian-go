# Architecture Decisions
This document outlines the architectural decisions and design patterns employed in the OpenSpec BIAN Go implementation. It describes the three-layer architecture, key components, and rationale behind technology choices.

## System Architecture

### Three-Layer Design

```
┌─────────────────────────────────────────────────────────────┐
│                    API Layer                                │
│  ┌─────────────────────┐  ┌─────────────────────────────┐   │
│  │   REST API          │  │   GraphQL API               │   │
│  │   /accounts/*       │  │   /graphql                  │   │
│  │   /transactions/*   │  │   /playground               │   │
│  │   /consents/*       │  │                             │   │
│  └─────────────────────┘  └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                 Domain Layer (BIAN Interfaces)             │
│  AccountService | TransactionService | BalanceService      │
│  ConsentService                                             │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                 Provider Layer                              │
│  Mock Provider (in-memory with realistic sample data)      │
│  Future: CDR, Plaid, Open Banking providers               │
└─────────────────────────────────────────────────────────────┘
```

**Key Principle**: Business logic lives in domain layer, transport is swappable.

### Dual API Architecture

Both REST and GraphQL APIs call domain interfaces directly (parallel implementations):
- **Zero coupling** between REST and GraphQL
- **No performance overhead** (both call domain directly)
- **Independent evolution** of each API
- **Provider implementations** swappable at runtime for both APIs

## 1. Domain Layer 

**Location:** `domains/`

Defines BIAN service interfaces following the [BIAN Service Landscape](https://bian.org/servicelandscape/) standard. Aligning to https://github.com/bian-official/public artefacts 

**Implemented Interfaces (BIAN v13.0.0 subset):**
- `AccountService` - Current Account Fulfillment service domain (read-only)
- `TransactionService` - Payment Execution service domain (read-only)
- `BalanceService` - Account Balance Management (read-only)
- `ConsentService` - Customer Consent Management (read-only)

**Future Interfaces:**
- `CustomerService` - Customer information, account associations
- Write operations (payment initiation, consent creation)

## 2. Model Layer

**Location:** `models/`

BIAN-aligned data structures representing banking entities.

**Core Types:**
- `Account` - Banking account with balance, status, metadata
- `Transaction` - Financial transaction with amount, timestamps, merchant info
- `Balance` - Available/current balance with currency
- `Money` - Decimal amount with currency (never float64)
- `Consent` - OAuth consent with scope, expiry, status

**Money Model Design Decision:**
- Uses `github.com/shopspring/decimal` for precise financial calculations
- JSON serialization to string format prevents JavaScript precision loss
- ISO 4217 currency codes (AUD, USD, GBP, EUR, CAD)
- Arithmetic operations: Add, Subtract, Multiply, Divide, Compare

## 3. REST API Layer

**Location:** `rest/`

RESTful HTTP API using Go 1.22+ `http.ServeMux` with zero external dependencies.

**Endpoints:**
- `GET /accounts/{id}` → RetrieveCurrentAccount
- `GET /accounts/{id}/balance` → RetrieveCurrentAccountBalance
- `GET /accounts/{id}/transactions` → RetrievePaymentTransactionHistory
- `GET /transactions/{id}` → RetrievePaymentTransaction
- `GET /consents/{id}` → RetrieveConsent
- `GET /consents/{id}/status` → RetrieveConsentStatus
- `GET /health` → Health check

**Middleware Stack:**
- Request logging (request ID, timing)
- Error recovery (panic handling)
- CORS handling
- Error response formatting

## 4. GraphQL API Layer

**Location:** `graphql/`

Modern API exposing BIAN domains via GraphQL with type-safe schema.

### Schema Structure

**File:** `graphql/schema.graphql`

**Query Operations:**
- `account(id)` - Account details
- `balance(accountId)` - Current balance
- `balances(accountId)` - All balance types
- `transaction(id)` - Transaction details
- `transactions(accountId, input)` - Transaction history with filtering
- `consent(id)` - Consent information
- `consentStatus(id)` - Consent status only

**Code Generation:**
- Uses [gqlgen](https://gqlgen.com/) for type-safe resolver generation
- Configuration in `gqlgen.yml`
- Run `go generate ./...` after schema changes
- Generated code in `graphql/generated/`

### Resolver Pattern

**File:** `graphql/resolver.go`

Root resolver wraps domain interfaces:
```go
Resolver {
    accountService
    transactionService
    balanceService
    consentService
}
```

Resolvers delegate to domain interfaces without business logic. All business rules live in domain/provider layers.

## 5. Provider Layer (Implementations)

**Location:** `providers/`

Region-specific implementations of domain interfaces.

**Current Providers:**
- `mock/` - In-memory provider with realistic test data (always available)

**Planned Providers:**
- `cdr/` - Australian Consumer Data Right implementation
- `plaid/` - US/Canada Plaid integration
- `obuk/` - UK Open Banking
- `psd2/` - EU PSD2

**Design Pattern:**
- Each provider implements all domain interfaces
- Normalization functions convert provider responses to canonical models
- OAuth token management abstracted
- Rate limiting and caching handled internally
- Configuration via struct (not global state)

### Mock Provider

**Sample Data:**
- 3 accounts (checking, savings, credit card)
- 11 transactions with realistic merchants
- 3 consents (active, expired, revoked)
- Multi-currency support (AUD, USD)

## 6. Unified Server

**Location:** `server/`

Combines REST and GraphQL APIs on single HTTP server:
- **Port**: Configurable via `PORT` environment variable (default: 8080)
- **Graceful shutdown**: 30-second timeout
- **Health checks**: `/health` endpoint
- **Development**: GraphQL Playground at `/playground`

## Technology Decisions

### Go 1.22+ Standard Library
- **http.ServeMux**: Enhanced with path variables, zero dependencies
- **context.Context**: Cancellation and timeout handling
- **encoding/json**: Standard JSON serialization

### External Dependencies (Minimal)
- **shopspring/decimal**: Precise financial arithmetic
- **99designs/gqlgen**: GraphQL server and code generation
- **google/uuid**: Request ID generation

### Deployment
- **Docker**: Multi-stage build with Alpine Linux
- **Kubernetes**: Stateless, horizontally scalable
- **Serverless**: Compatible with AWS Lambda, Google Cloud Run
- **Health checks**: Built-in `/health` endpoint
