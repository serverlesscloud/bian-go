# BIAN Go

> Banking Industry Architecture Network interfaces for Go

A Go library providing BIAN-compliant banking interfaces with dual API layer (REST + GraphQL). Built for developers creating fintech applications with open banking integration.

**Repository:** `github.com/serverlesscloud/bian-go`

## 🎆 Key Features

- **🌍 Cloud-portable** - Deploy anywhere (Docker, Kubernetes, serverless)
- **🔄 Provider-agnostic** - Swap between CDR, Plaid, Open Banking implementations
- **🚀 Dual API** - Both REST and GraphQL on single server
- **🔒 Type-safe** - Compile-time guarantees with Go's type system
- **💰 Decimal precision** - No floating-point errors in financial calculations
- **🏦 BIAN v13.0.0 aligned** - Standardized banking operations
- **🧪 Mock provider** - Development without external dependencies
- **⚙️ Zero dependencies** - REST API uses only Go standard library

## 🚀 Quick Start

```bash
go get github.com/serverlesscloud/bian-go
```

```go
import (
    "github.com/serverlesscloud/bian-go/providers/mock"
    "github.com/serverlesscloud/bian-go/server"
)

provider := mock.NewProvider()
config := server.DefaultConfig()
srv := server.NewServer(provider, provider, provider, provider, config)
srv.Start() // Starts both REST and GraphQL APIs
```

**Server URLs:**
- REST API: http://localhost:8080/
- GraphQL API: http://localhost:8080/graphql
- GraphQL Playground: http://localhost:8080/playground
- Health Check: http://localhost:8080/health

See `examples/` directory for complete working implementations.

## 🏢 Architecture

### Three-Layer Design

```
API Layer (REST + GraphQL)
    ↓
Domain Layer (BIAN Interfaces)
    ↓
Provider Layer (CDR/Plaid/Open Banking)
```

**Key principle:** Business logic lives in domain layer, transport is swappable.

### Dual API Support

Both REST and GraphQL APIs call domain interfaces directly (parallel implementations):
- **Zero coupling** between REST and GraphQL
- **No performance overhead** (both call domain directly)
- **Independent evolution** of each API

### Supported Standards

| Region | Standard | Provider Package | Status |
|--------|----------|------------------|--------|
| 🇦🇺 Australia | Consumer Data Right (CDR) | `providers/cdr` | Planned |
| 🇺🇸 🇨🇦 US/Canada | Plaid | `providers/plaid` | Planned |
| 🇬🇧 UK | Open Banking | `providers/obuk` | Planned |
| 🇪🇺 EU | PSD2 | `providers/psd2` | Planned |
| 🧪 Testing | Mock | `providers/mock` | **Available** |

## 📁 Project Structure

```
bian-go/
├── domains/              # BIAN service interfaces
│   ├── accounts.go
│   ├── transactions.go
│   ├── consents.go
│   └── balance.go
│
├── models/               # Canonical data types
│   ├── account.go
│   ├── transaction.go
│   ├── money.go
│   └── enums.go
│
├── rest/                 # REST API layer
│   ├── handlers.go
│   ├── middleware.go
│   └── server.go
│
├── graphql/              # GraphQL API layer
│   ├── schema.graphql
│   ├── resolver.go
│   ├── server.go
│   └── generated/        # gqlgen output
│
├── providers/            # Banking implementations
│   └── mock/            # Testing provider
│
├── server/               # Unified server
└── examples/            # Working examples
```

## 🔌 Core Interfaces

See `domains/` for complete interface definitions. Key patterns:

- **AccountService** - Account lifecycle and balance queries (BIAN Current Account Fulfillment)
- **TransactionService** - Payment execution and transaction history (BIAN Payment Execution)
- **BalanceService** - Balance information (BIAN Account Balance Management)
- **ConsentService** - OAuth consent management (BIAN Customer Consent Management)

All interfaces accept `context.Context` as first parameter for cancellation/timeouts.

## 🌐 REST API

**Base URL:** http://localhost:8080/

### Account Endpoints
```bash
# Get account details
GET /accounts/{id}

# Get current balance
GET /accounts/{id}/balance

# Get all balance types
GET /accounts/{id}/balances

# Get transaction history
GET /accounts/{id}/transactions?fromDate=2024-01-01&limit=10
```

### Transaction Endpoints
```bash
# Get transaction details
GET /transactions/{id}
```

### Consent Endpoints
```bash
# Get consent details
GET /consents/{id}

# Get consent status
GET /consents/{id}/status
```

## 🔍 GraphQL API

**Endpoint:** http://localhost:8080/graphql  
**Playground:** http://localhost:8080/playground

### Example Queries

```graphql
# Get account with balance
query {
  account(id: "acc-001") {
    id
    accountNumber
    accountType
    productName
    status
  }
  
  balance(accountId: "acc-001") {
    balanceType
    amount {
      amount
      currency
    }
  }
}

# Get transaction history
query {
  transactions(accountId: "acc-001", input: {
    limit: 10
    fromDate: "2024-01-01"
  }) {
    id
    transactionType
    amount {
      amount
      currency
    }
    description
    postingDate
  }
}
```

Configuration via `gqlgen.yml`. Run `go generate ./...` to regenerate after schema changes.

## 📦 Provider Implementation

Implement domain interfaces to create new providers:

```go
type MyProvider struct {
    client *BankingAPIClient
}

func (p *MyProvider) RetrieveCurrentAccount(ctx context.Context, id string) (*models.Account, error) {
    // 1. Call external banking API
    // 2. Normalize to BIAN canonical model
    // 3. Return standardized Account
}
```

See `providers/mock/` for reference implementation.

## 💰 Money Model

Precise decimal arithmetic for financial calculations:

```go
// Create money amounts
aud, _ := models.NewMoney(decimal.NewFromFloat(123.45), "AUD")
usd, _ := models.NewMoney(decimal.NewFromFloat(100.00), "USD")

// Arithmetic operations
sum, _ := aud.Add(otherAUD)
product := aud.Multiply(decimal.NewFromFloat(1.5))

// Comparisons
equal := aud.Equal(otherAUD)
greater, _ := aud.GreaterThan(otherAUD)

// JSON serialization (string format for precision)
// {"amount": "123.45", "currency": "AUD"}
```

**Supported Currencies:** AUD, USD, GBP, EUR, CAD, JPY, CHF, CNY, SEK, NZD

## 🧪 Mock Provider

Includes realistic sample data for development:

### Sample Accounts
- `acc-001`: Checking Account (AUD $2,547.83)
- `acc-002`: Savings Account (AUD $15,420.50)
- `acc-003`: Credit Card (USD -$1,250.75)

### Sample Transactions
- 11 transactions across accounts
- Various types: debit, credit, transfer, payment, fee
- Realistic merchants: Woolworths, Energy Australia, Amazon

### Sample Consents
- `consent-001`: Active (account:read, transaction:read, balance:read)
- `consent-002`: Expired
- `consent-003`: Revoked

## 🔧 Development

```bash
# Install dependencies
go mod download

# Generate GraphQL code
go generate ./...

# Run tests
go test ./...

# Run example server
cd examples/basic && go run main.go

# Build Docker image
docker build -t bian-go .

# Run in Docker
docker run -p 8080:8080 bian-go
```

### Environment Variables
- `PORT`: Server port (default: 8080)
- `ENABLE_PLAYGROUND`: Enable GraphQL Playground (default: true)

## 🔄 BIAN Spec Synchronization

Automated workflow syncs with BIAN releases:
- Weekly checks for new Service Landscape versions
- Auto-generates PRs with updated OpenAPI specs
- Validates breaking changes

See `.github/workflows/bian-sync.yml`

## 📝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

Key areas for contribution:
- New provider implementations (regional banking standards)
- Performance optimizations
- Documentation improvements
- Test coverage expansion

## 👥 Contributors

- **Amazon Q Developer** - 3 Musketeers build system implementation, DevOps automation

## 📋 License

Apache License 2.0 - See [LICENSE](LICENSE)

Patent protection, corporate-friendly, permissive for commercial use.