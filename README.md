# BIAN Go

> Banking Industry Architecture Network interfaces for Go

A Go library providing BIAN-compliant banking interfaces with GraphQL API layer. Built for developers creating fintech applications with open banking integration.

**Repository:** `github.com/serverlesscloud/bian-go`

## Overview

BIAN-Go provides standardized interfaces for banking operations aligned with the [BIAN Service Landscape](https://bian.org/servicelandscape/). The library enables:

- **Cloud-portable** - Deploy anywhere (Docker, Kubernetes, cloud platforms)
- **Provider-agnostic** - Swap between CDR, Plaid, Open Banking implementations
- **GraphQL-first** - Modern API with real-time subscriptions
- **Type-safe** - Compile-time guarantees with Go's type system

### Core Philosophy

1. **Interfaces over implementations** - BIAN domains define contracts
2. **Portability over convenience** - No cloud vendor lock-in
3. **Standards over opinions** - Align with BIAN specifications exactly
4. **Testability over production-only code** - Mock providers included

## Quick Start

```bash
go get github.com/serverlesscloud/bian-go
```

```go
import (
    "github.com/serverlesscloud/bian-go/domains"
    "github.com/serverlesscloud/bian-go/graphql"
    "github.com/serverlesscloud/bian-go/providers/mock"
)

provider := mock.NewProvider()
server := graphql.NewServer(provider)
http.ListenAndServe(":8080", server)
```

See `examples/` directory for complete working implementations.

## Architecture

### Three-Layer Design

```
API Layer (GraphQL)
    ↓
Domain Layer (BIAN Interfaces)
    ↓
Provider Layer (CDR/Plaid/Open Banking)
```

**Key principle:** Business logic lives in domain layer, transport is swappable.

### Supported Standards

| Region | Standard | Provider Package |
|--------|----------|------------------|
| 🇦🇺 Australia | Consumer Data Right (CDR) | `providers/cdr` |
| 🇺🇸 🇨🇦 US/Canada | Plaid | `providers/plaid` |
| 🇬🇧 UK | Open Banking | `providers/obuk` (planned) |
| 🇪🇺 EU | PSD2 | `providers/psd2` (planned) |
| 🧪 Testing | Mock | `providers/mock` |

## Project Structure

```
bian-go/
├── domains/              # BIAN service interfaces
│   ├── accounts.go
│   ├── transactions.go
│   ├── consents.go
│   └── customers.go
│
├── models/               # Canonical data types
│   ├── account.go
│   ├── transaction.go
│   ├── money.go
│   └── enums.go
│
├── graphql/              # GraphQL API layer
│   ├── schema.graphql
│   ├── resolver.go
│   ├── server.go
│   └── generated/        # gqlgen output
│
├── providers/            # Banking implementations
│   ├── mock/            # Testing provider
│   ├── cdr/             # Australian CDR
│   └── plaid/           # US/Canada
│
├── normalizers/          # Data transformation
├── validators/           # Business rules
└── examples/            # Working examples
```

## Core Interfaces

See `domains/` for complete interface definitions. Key patterns:

- **AccountService** - Account lifecycle and balance queries
- **TransactionService** - Payment execution and transaction history
- **ConsentService** - OAuth consent management
- **CustomerService** - Customer information

All interfaces accept `context.Context` as first parameter for cancellation/timeouts.

## GraphQL API

Schema located in `graphql/schema.graphql`. Supports:

- **Queries** - Account/transaction retrieval
- **Mutations** - Payment initiation, consent management
- **Subscriptions** - Real-time balance/transaction updates via WebSocket

Configuration via `gqlgen.yml`. Run `go generate ./...` to regenerate after schema changes.

## Provider Implementation

Implement domain interfaces to create new providers:

```go
type MyProvider struct {
    client *BankingAPIClient
}

func (p *MyProvider) GetAccount(ctx context.Context, id string) (*models.Account, error) {
    // 1. Call external banking API
    // 2. Normalize to BIAN canonical model
    // 3. Return standardized Account
}
```

See `providers/mock/` for reference implementation.

## BIAN Spec Synchronization

Automated workflow syncs with BIAN releases:
- Weekly checks for new Service Landscape versions
- Auto-generates PRs with updated OpenAPI specs
- Validates breaking changes

See `.github/workflows/bian-sync.yml`

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

Key areas for contribution:
- New provider implementations (regional banking standards)
- Performance optimizations
- Documentation improvements
- Test coverage expansion

## License

Apache License 2.0 - See [LICENSE](LICENSE)

Patent protection, corporate-friendly, permissive for commercial use.