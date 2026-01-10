# Design: MVP Foundation Architecture

## Context

BIAN-Go implements the Banking Industry Architecture Network (BIAN) standard in Go with dual API layers (REST and GraphQL). The project must support:
- Cloud-portable deployment (Docker, Kubernetes, serverless)
- Provider-agnostic implementation (swap between CDR, Plaid, Open Banking)
- Type-safe interfaces with Go's type system
- BIAN v13 Service Landscape alignment

**Key Constraints:**
- No existing codebase - greenfield implementation
- Must work without external banking API dependencies (mock provider)
- Follow idiomatic Go coding style
- Containerized deployment model

**Stakeholders:**
- Fintech application developers using BIAN-Go library
- Banking integration teams implementing provider adapters
- BIAN specification maintainers (alignment requirement)

## Goals / Non-Goals

### Goals
- Establish three-layer architecture (API Layer → Domain → Provider)
- Implement dual API support (REST + GraphQL) with parallel implementations
- Implement essential BIAN service domains (Account, Transaction, Balance, Consent)
- Create foundation models with proper Money handling (decimal-based)
- Provide working mock provider for development/testing
- Enable read-only REST and GraphQL queries
- Maintain strict BIAN naming and structure alignment

### Non-Goals
- Write operations (payment initiation, consent creation) - deferred to future changes
- Real-time subscriptions via WebSocket - deferred
- External provider implementations (CDR, Plaid) - deferred
- Full BIAN operation coverage - MVP uses essential subset only
- Performance optimization - premature at this stage
- Authentication/authorization - out of scope for core library

## Decisions

### Decision 1: Subset BIAN Operations Strategy

**Choice:** Implement essential read-only BIAN operations first, add incrementally

**Rationale:**
- BIAN Current Account domain defines 30+ operations (InitiateCurrentAccountFacility, UpdateCurrentAccountFacility, RegisterDirectDebit, etc.)
- MVP needs only core read operations: RetrieveCurrentAccount, RetrieveCurrentAccountBalance
- Implemented operations follow BIAN specification exactly (naming, parameters, semantics)
- Future additions simply add more operations to existing service domains
- Reduces implementation scope from ~150 operations to ~10 for MVP

**Alternatives Considered:**
1. **Full BIAN implementation** - Comprehensive but 10-15x effort, delays MVP by months
2. **Custom operation names** - Faster but loses BIAN compliance, migration pain later
3. **Hybrid approach** - Mix BIAN and custom - creates confusion, specification drift

**Trade-offs:**
- ✅ Faster time to working MVP (weeks vs months)
- ✅ Maintains BIAN compliance for implemented operations
- ✅ Clear upgrade path (add operations incrementally)
- ❌ Incomplete BIAN service domain coverage initially
- ❌ Documentation must clearly state "subset" status

### Decision 2: Decimal-Based Money Model

**Choice:** Use `github.com/shopspring/decimal` for Money type with ISO 4217 currencies

**Rationale:**
- Financial calculations require exact decimal arithmetic (no floating-point errors)
- Industry standard: `0.1 + 0.2 = 0.3` (not `0.30000000000000004`)
- shopspring/decimal is battle-tested, widely used in Go financial applications
- Supports all required operations: add, subtract, multiply, divide, compare
- JSON serialization to string format preserves precision
- Currency stored as 3-letter ISO 4217 code (AUD, USD, GBP, EUR, CAD)

**Alternatives Considered:**
1. **Integer cents/pence** - Simple but currency-specific conversion logic, loses semantic clarity
2. **Third-party money library** (github.com/Rhymond/go-money) - Adds opinion on currency handling we don't need
3. **float64** - ❌ NEVER for financial data (rounding errors cause real money loss)

**Trade-offs:**
- ✅ Prevents floating-point precision errors
- ✅ JSON string format prevents JavaScript number precision loss
- ✅ Idiomatic Go struct with value semantics
- ➖ Adds external dependency (but well-maintained, stable)
- ➖ Slightly slower than integer arithmetic (acceptable for MVP)

**Example:**
```go
type Money struct {
    Amount   decimal.Decimal `json:"amount"`   // e.g., "123.45"
    Currency string          `json:"currency"` // e.g., "AUD"
}

// JSON representation:
// {"amount": "123.45", "currency": "AUD"}
```

### Decision 3: Dual API Architecture with Parallel Implementations

**Choice:** Both REST and GraphQL APIs call domain interfaces directly (parallel implementations)

**Rationale:**
- Aligns with project philosophy: "Interfaces over implementations"
- Domain layer defines BIAN service contracts (Go interfaces)
- Both API layers are independent transport mechanisms (swappable)
- REST and GraphQL resolvers/handlers are thin adapters (no business logic)
- Provider layer has regional implementations (mock, CDR, Plaid)
- Business logic lives in domain interfaces, not in API handlers
- Enables testing with mock provider without external dependencies
- Parallel approach avoids coupling and latency overhead (vs GraphQL-wraps-REST pattern)

**Architecture:**
```
┌──────────────────────────┐  ┌──────────────────────────┐
│   REST API Layer         │  │   GraphQL API Layer      │
│  (handlers, middleware)  │  │  (schema, resolvers)     │
│  GET /accounts/{id}      │  │  query { account(id) }   │
│  GET /transactions       │  │  query { transactions }  │
│  - Input validation      │  │  - Input validation      │
│  - Error mapping         │  │  - Error mapping         │
└──────────────────────────┘  └──────────────────────────┘
              ↓ calls                   ↓ calls
              └───────────┬─────────────┘
                          ↓
              ┌─────────────────────────────────────┐
              │     Domain Layer (Interfaces)       │
              │  - AccountService                   │
              │  - TransactionService               │
              │  - BalanceService                   │
              │  - ConsentService                   │
              └─────────────────────────────────────┘
                          ↓ implements
              ┌─────────────────────────────────────┐
              │     Provider Layer                  │
              │  - mock.Provider (in-memory)        │
              │  - cdr.Provider (future)            │
              │  - plaid.Provider (future)          │
              └─────────────────────────────────────┘
```

**Alternatives Considered:**
1. **GraphQL wraps REST** - GraphQL resolvers make HTTP calls to internal REST endpoints
   - ❌ Adds latency (internal HTTP round-trip)
   - ❌ Couples GraphQL to REST implementation
   - ❌ Complicates error handling (HTTP errors → GraphQL errors)
   - ❌ Harder to test (need HTTP server running)

2. **REST wraps GraphQL** - REST endpoints execute GraphQL queries internally
   - ❌ Unusual pattern, adds unnecessary GraphQL query parsing overhead
   - ❌ Couples REST to GraphQL schema
   - ❌ Makes REST dependent on GraphQL server availability

3. **Shared controller layer** - Both APIs call a controller/service layer above domain
   - ❌ Adds unnecessary abstraction (domain interfaces already serve this role)
   - ❌ More boilerplate code
   - ➖ Only beneficial if significant business logic exists outside domain (not our case)

4. **Monolithic unified API** - Single codebase that outputs both REST and GraphQL
   - ❌ Tightly couples both API formats
   - ❌ Harder to maintain and reason about
   - ❌ Difficult to optimize each API independently

**Trade-offs:**
- ✅ Zero coupling between REST and GraphQL (independent evolution)
- ✅ No performance overhead (both call domain directly)
- ✅ Easy to test (mock domain interfaces for each API independently)
- ✅ Provider implementations swappable at runtime for both APIs
- ✅ Either API can be disabled/removed without affecting the other
- ✅ Clear separation of concerns (REST adapts HTTP, GraphQL adapts schema)
- ➖ Slight code duplication (validation, error handling in both APIs)
- ➖ Two sets of integration tests (REST tests, GraphQL tests)
- ✅ Duplication is minimal and worth independence gained

### Decision 4: Mock Provider Implementation Scope

**Choice:** In-memory mock provider with realistic BIAN-aligned test data

**Rationale:**
- Enables development/testing without external banking API accounts
- Provides realistic data for UI/UX testing
- Demonstrates correct provider implementation pattern
- Serves as reference for external provider implementations
- Must follow BIAN data structures exactly (not simplified)

**Mock Data Scope:**
- 3-5 sample accounts (checking, savings, credit card)
- 20-30 sample transactions per account (various merchants, amounts)
- 2-3 sample consents (active, expired)
- Multi-currency support (AUD, USD, GBP)
- Realistic timestamps, balances, transaction types

**Alternatives Considered:**
1. **Minimal mock (1 account, 3 transactions)** - Too sparse for realistic testing
2. **Generated random data** - Non-deterministic, hard to test
3. **External test API** - Adds dependency, increases complexity

**Trade-offs:**
- ✅ Zero external dependencies
- ✅ Deterministic test data (reproducible tests)
- ✅ Fast execution (in-memory)
- ❌ Not suitable for production (obviously)
- ❌ Doesn't test external API integration (that's for provider-specific tests)

### Decision 5: GraphQL Code Generation

**Choice:** Use `github.com/99designs/gqlgen` for type-safe GraphQL server

**Rationale:**
- Industry standard for GraphQL in Go (most popular, active maintenance)
- Generates type-safe resolver interfaces from schema
- Schema-first approach (define schema.graphql, generate Go code)
- Supports custom scalar types (Money, DateTime)
- Configuration via gqlgen.yml
- `go generate ./...` workflow integrates with standard Go tooling

**Configuration:**
- Schema location: `graphql/schema.graphql`
- Generated code: `graphql/generated/`
- Custom resolvers: `graphql/resolver.go`
- Model mapping: domain models → GraphQL types

**Alternatives Considered:**
1. **graphql-go/graphql** - Older, code-first approach (less type-safe)
2. **thunder (Samsara)** - Opinionated, limited adoption
3. **Manual GraphQL handling** - Error-prone, defeats purpose of schema

**Trade-offs:**
- ✅ Type safety at compile time
- ✅ Schema-first ensures API contract clarity
- ✅ Auto-generates boilerplate
- ➖ Adds build step (`go generate`)
- ➖ Learning curve for gqlgen-specific patterns

### Decision 6: REST API Design with Standard Library

**Choice:** Use Go 1.22+ `http.ServeMux` with standard RESTful conventions (not BIAN-style operations)

**Rationale:**
- Go 1.22 enhanced ServeMux supports path variables (e.g., `/accounts/{id}`) - zero external dependencies
- Standard RESTful conventions are intuitive for HTTP clients (GET /accounts/{id}, not POST /accounts/{id}/retrieve)
- BIAN operation semantics map cleanly to HTTP verbs:
  - `RetrieveCurrentAccount` → `GET /accounts/{id}`
  - `RetrievePaymentTransactionHistory` → `GET /accounts/{id}/transactions`
  - `RetrieveConsent` → `GET /consents/{id}`
- RESTful approach aligns with OpenAPI/Swagger standards for documentation
- Middleware pattern for cross-cutting concerns (logging, error handling, CORS)
- JSON responses use canonical domain models (same as GraphQL)

**REST Endpoint Design:**
```
Account Service:
  GET /accounts/{id}              → RetrieveCurrentAccount
  GET /accounts/{id}/balance      → RetrieveCurrentAccountBalance

Transaction Service:
  GET /transactions/{id}          → RetrievePaymentTransaction
  GET /accounts/{id}/transactions → RetrievePaymentTransactionHistory
    Query params: fromDate, toDate, limit, offset

Consent Service:
  GET /consents/{id}              → RetrieveConsent
  GET /consents/{id}/status       → RetrieveConsentStatus
```

**Alternatives Considered:**
1. **BIAN-style REST operations** - POST `/accounts/{id}/retrieve`
   - ✅ More BIAN-compliant naming
   - ❌ Violates HTTP semantics (GET for retrieval, not POST)
   - ❌ Less intuitive for REST clients
   - ❌ Doesn't align with OpenAPI best practices

2. **Chi router** - Lightweight third-party router
   - ✅ Composable middleware, good developer experience
   - ➖ Adds external dependency (albeit small)
   - ➖ Go 1.22+ ServeMux now provides similar functionality

3. **Gin framework** - Full-featured HTTP framework
   - ✅ Rich features (validation, binding, rendering)
   - ❌ Opinionated, larger dependency footprint
   - ❌ Overkill for simple CRUD operations
   - ❌ Adds learning curve for contributors

4. **Gorilla Mux** - Mature routing library
   - ❌ Now in maintenance mode (project archived)
   - ❌ Use Chi or standard library instead

**Middleware Stack:**
- Request logging (request ID, method, path, duration)
- Error recovery (panic handling)
- CORS handling (configurable origins)
- Error response formatting (consistent JSON error structure)

**Error Response Format:**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Account not found",
    "details": "No account exists with ID acc-99999"
  }
}
```

**Trade-offs:**
- ✅ Zero external dependencies (standard library only)
- ✅ Familiar HTTP semantics for REST clients
- ✅ Easy to document with OpenAPI/Swagger
- ✅ Lightweight and performant
- ✅ Simple middleware pattern for cross-cutting concerns
- ➖ Slightly less BIAN-aligned naming (but semantics preserved)
- ➖ Manual request/response handling (vs framework's binding/validation)
- ✅ Manual handling gives full control and transparency

## Risks / Trade-offs

### Risk 1: BIAN Specification Drift
**Risk:** BIAN updates service domains, our subset becomes outdated

**Mitigation:**
- Document BIAN version alignment in `.bian-version` file (start with v13.0.0)
- Automated workflow for BIAN spec updates (planned in `project.md`)
- Clear documentation of implemented vs. deferred operations
- Version compatibility matrix in README

**Impact:** Medium - requires ongoing maintenance

### Risk 2: Decimal Dependency Maintenance
**Risk:** shopspring/decimal library becomes unmaintained

**Mitigation:**
- Library has 4000+ GitHub stars, wide adoption, active maintenance
- Isolated to `models/money.go` - single file to swap if needed
- Standard decimal arithmetic - easy to replace with alternative
- Consider: abstract Money interface if risk materializes

**Impact:** Low - library is stable and widely used

### Risk 3: Mock Provider Divergence
**Risk:** Mock provider behavior diverges from real provider implementations

**Mitigation:**
- Mock provider must implement same domain interfaces (enforced by Go compiler)
- Integration tests for mock provider document expected behavior
- When real providers added (CDR, Plaid), validate against same test suite
- Mock data follows BIAN canonical models exactly

**Impact:** Low - interface contracts prevent divergence

### Risk 4: API Evolution and Breaking Changes
**Risk:** Future changes to REST or GraphQL APIs break existing clients

**Mitigation:**
- **GraphQL**: Follow best practices for schema evolution (additive changes only), use deprecation annotations
- **REST**: Version API via URL path (e.g., `/v1/accounts/{id}`) or Accept header if breaking changes needed
- Document all API changes in CHANGELOG with migration guides
- Consider API versioning strategy early (v1 can be implicit for MVP)
- Use OpenAPI spec for REST API documentation and contract testing

**Impact:** Medium - requires discipline in API design and version management

### Risk 5: REST and GraphQL Consistency
**Risk:** REST and GraphQL APIs diverge in behavior or data representation

**Mitigation:**
- Both APIs use same canonical domain models (enforced by Go type system)
- Both APIs call same domain interface methods (single source of truth)
- Shared test fixtures for integration tests (verify same data from both APIs)
- Document API parity expectations in README
- Consider contract tests that verify REST and GraphQL return equivalent data

**Impact:** Low - architecture enforces consistency at domain layer

## Migration Plan

**N/A** - This is the initial implementation with no existing code to migrate.

**For Future Provider Implementations:**
1. Implement domain interfaces for new provider (e.g., `cdr.Provider`)
2. Write provider-specific integration tests
3. Add normalization functions (external API → canonical models)
4. Document provider-specific configuration (OAuth, rate limits)
5. Update README with provider status matrix

**For Future Operation Additions:**
1. Update domain interface (add new method)
2. Implement in mock provider first (for testing)
3. Add REST endpoint handler
4. Update GraphQL schema (add mutation/query)
5. Regenerate GraphQL code (`go generate ./...`)
6. Update integration tests (REST and GraphQL)
7. Document in CHANGELOG

## Open Questions

1. **Error Handling Strategy**
   - How should domain interfaces communicate errors? (Go error interface, custom error types, error codes?)
   - Proposal: Start with standard Go errors, introduce typed errors when patterns emerge
   - Decision: Defer to implementation - use Go errors initially

2. **Context Propagation**
   - What should we pass in context.Context? (request ID, user ID, correlation ID?)
   - Proposal: Minimal context for MVP (just cancellation), extend as needed
   - Decision: Defer to implementation - context.Context for cancellation only

3. **Rate Limiting**
   - Should mock provider simulate rate limiting for realistic testing?
   - Proposal: No rate limiting in mock (deterministic), add to real providers only
   - Decision: Defer to real provider implementations

4. **BIAN Control Records**
   - Should we model BIAN Control Records explicitly or embed in domain models?
   - Proposal: Embed in canonical models for MVP (e.g., `Account.ID` is Control Record ID)
   - Decision: Document in code comments, consider explicit modeling in future change

5. **Multi-tenancy**
   - How should the library handle multi-tenant scenarios?
   - Proposal: Out of scope for MVP - consumer responsibility (pass tenant ID in context)
   - Decision: Document as consumer responsibility, revisit if common pattern emerges
