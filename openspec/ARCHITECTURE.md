# Architecture Decisions
This document outlines the architectural decisions and design patterns employed in the OpenSpec BIAN Go implementation. It describes the three-layer architecture, key components, and rationale behind technology choices.

#### 1 Domain Layer 
Defines BIAN service interfaces following the [BIAN Service Landscape](https://bian.org/servicelandscape/) standard. Aligning to https://github.com/bian-official/public artefacts 

**Sample intefaces**
- `AccountService` - Account retrieval, balance queries, account history
- `TransactionService` - Transaction listing, payment initiation, status tracking
- `ConsentService` - OAuth consent creation, revocation, refresh
- `CustomerService` - Customer information, account associations

#### 2 Model layer**

BIAN-aligned data structures representing banking entities.

**Core Types:**
- `Account` - Banking account with balance, status, metadata
- `Transaction` - Financial transaction with amount, timestamps, merchant info
- `Balance` - Available/current balance with currency
- `Money` - Decimal amount with currency (never float64)
- `Consent` - OAuth consent with scope, expiry, status

### 4. Provider Layer (Implementations)

**Location:** `providers/`

Region-specific implementations of domain interfaces.

**Providers:**
- `mock/` - In-memory provider with realistic test data (always available)
- `cdr/` - Australian Consumer Data Right implementation
- `plaid/` - US/Canada Plaid integration
- `obuk/` - UK Open Banking (planned)
- `psd2/` - EU PSD2 (planned)

**Design Pattern:**
- Each provider implements all domain interfaces
- Normalization functions convert provider responses to canonical models
- OAuth token management abstracted
- Rate limiting and caching handled internally
- Configuration via struct (not global state)


## GraphQL API Layer

**Location:** `graphql/`

Modern API exposing BIAN domains via GraphQL with real-time capabilities.

### Schema Structure

**File:** `graphql/schema.graphql`

Three operation types:
1. **Query** - Read operations (accounts, transactions, balance)
2. **Mutation** - Write operations (payment initiation, consent management)
3. **Subscription** - Real-time updates (balance changes, new transactions)

**Code Generation:**
- Uses [gqlgen](https://gqlgen.com/) for type-safe resolver generation
- Configuration in `gqlgen.yml`
- Run `go generate ./...` after schema changes
- Generated code in `graphql/generated/`

### Resolver Pattern

**File:** `graphql/resolver.go`

Root resolver wraps domain interfaces:
```
Resolver {
    accountService
    transactionService
    consentService
    customerService
    pubsub (for subscriptions)
}
```

Resolvers delegate to domain interfaces without business logic. All business rules live in domain/provider layers.
