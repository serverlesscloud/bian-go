# Project Context

## Purpose
BIAN-Go is an open-source Go library implementing the Banking Industry Architecture Network (BIAN) standard with dual API layers (REST + GraphQL). The library enables cloud-portable fintech application development with standardized banking integration across multiple regions (Australia CDR, US/Canada Plaid, UK Open Banking).

## Key Features Implemented

### 🏦 BIAN v13.0.0 Compliance
- **AccountService** - Current Account Fulfillment service domain (read-only subset)
- **TransactionService** - Payment Execution service domain (read-only subset)
- **BalanceService** - Account Balance Management (read-only subset)
- **ConsentService** - Customer Consent Management (read-only subset)

### 🚀 Dual API Architecture
- **REST API** - RESTful endpoints using Go 1.22+ http.ServeMux (zero dependencies)
- **GraphQL API** - Type-safe schema with gqlgen, interactive playground
- **Parallel Implementation** - Both APIs call domain interfaces directly (no coupling)
- **Unified Server** - Single HTTP server hosting both APIs

### 💰 Precise Financial Calculations
- **Decimal Money Model** - Uses shopspring/decimal (no floating-point errors)
- **Multi-currency Support** - ISO 4217 codes (AUD, USD, GBP, EUR, CAD, etc.)
- **Arithmetic Operations** - Add, Subtract, Multiply, Divide, Compare
- **JSON Serialization** - String format prevents JavaScript precision loss

### 🧪 Mock Provider
- **In-memory Implementation** - No external dependencies required
- **Realistic Sample Data** - 3 accounts, 11 transactions, 3 consents
- **Multi-currency Accounts** - AUD and USD examples
- **Various Transaction Types** - Debit, credit, transfer, payment, fee

### ⚙️ Production Ready
- **Middleware Stack** - Request logging, error recovery, CORS, panic handling
- **Graceful Shutdown** - 30-second timeout with signal handling
- **Health Checks** - `/health` endpoint for load balancers
- **Docker Support** - Multi-stage build with Alpine Linux
- **Environment Configuration** - Port, playground enable/disable

### 📊 Testing & Quality
- **Unit Tests** - Money model with 100% coverage
- **Integration Tests** - REST and GraphQL API testing
- **Type Safety** - Go compiler enforces interface contracts
- **Error Handling** - Consistent error responses across APIs

## Tech Stack
- **Go 1.22+** - Enhanced http.ServeMux with path variables
- **GraphQL** - 99designs/gqlgen for type-safe code generation
- **Decimal Arithmetic** - shopspring/decimal for financial precision
- **Docker** - Multi-stage containerization
- **OpenSpec** - Architecture and change management

## Project Conventions

### Code Style
Follow idiomatic go coding style guide here https://google.github.io/styleguide/go/

### Architecture Patterns
See `ARCHITECTURE.md`

### Testing Strategy
See `TESTING.md`

### Git Workflow
- Trunk based development with tagged releases pegged against BIAN

### BIAN Specification Synchronization

#### Automated Workflow

**File:** `.github/workflows/bian-sync.yml`

**Schedule:** Weekly (Sundays 00:00 UTC)

**Process:**
1. Check BIAN website for latest Service Landscape version
2. Download OpenAPI specifications for relevant domains
3. Generate Go code using `oapi-codegen`
4. Run tests to detect breaking changes
5. Create PR with updates and changelog
6. Notify maintainers for review

**Manual Trigger:**
- Workflow can be manually triggered via GitHub Actions UI
- Specify BIAN version if not latest

#### Generated Code Locations

- `domains/bian/{domain}/types.go` - Generated from OpenAPI schemas
- `domains/bian/{domain}/client.go` - Generated client interfaces
- Version tracked in `.bian-version` file (currently v13.0.0)

## Domain Context
- BIAN https://bian.org/servicelandscape/

## Important Constraints
The code base should be deployabled via containerised artefacts to run either locally, in serverless runtims or long lived Kubernetes clusters.

## External Dependencies
- **shopspring/decimal** - Precise financial arithmetic
- **99designs/gqlgen** - GraphQL server and code generation
- **google/uuid** - Request ID generation

## Deployment Options

### Local Development
```bash
cd examples/basic && go run main.go
```

### Docker
```bash
docker build -t bian-go .
docker run -p 8080:8080 bian-go
```

### Kubernetes
- Stateless application (no persistent storage)
- Horizontal scaling supported
- Health check endpoint: `/health`
- Configurable via environment variables

### Serverless
- Compatible with AWS Lambda, Google Cloud Run
- Fast cold start (minimal dependencies)
- Stateless design
