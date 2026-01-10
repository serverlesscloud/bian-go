# Project Context

## Purpose
BIAN-Go is an open-source Go library implementing the Banking Industry Architecture Network (BIAN) standard with a GraphQL API layer. The library enables cloud-portable fintech application development with standardized banking integration across multiple regions (Australia CDR, US/Canada Plaid, UK Open Banking).

## Tech Stack
- Go
- GraphQL
- OpenSpec

## Project Conventions

### Code Style
Follow idiomatic go coding style guide here https://google.github.io/styleguide/go/

### Architecture Patterns
See `ARCHITECTURE.md`

### Testing Strategy

#### Unit Tests

**Domain Layer:**
- Test interfaces with mock implementations
- Verify error handling
- Test edge cases (nil values, empty strings)

**Model Layer:**
- Test JSON serialization/deserialization
- Verify enum validation
- Test Money arithmetic (if any)

**GraphQL Layer:**
- Test resolvers with mock providers
- Verify query/mutation execution
- Test subscription channel behavior

#### Integration Tests

**Provider Tests:**
- Use `httptest` to mock external APIs
- Test full normalization pipeline
- Verify OAuth token refresh
- Test rate limiting behavior

**GraphQL Tests:**
- Execute actual GraphQL queries against test server
- Verify complete request/response cycle
- Test error responses
- Test subscription setup/teardown

### Mock Provider

The `mock` provider serves dual purpose:
1. **Development** - Run server without external dependencies
2. **Testing** - Predictable test data for CI/CD

Includes realistic sample data:
- Multiple account types (checking, savings, investment)
- Transaction history with various merchants
- Active and expired consents
- Customer information

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
- Version tracked in `.bian-version` file


## Domain Context
- BIAN https://bian.org/servicelandscape/

## Important Constraints
The code base should be deployabled via containerised artefacts to run either locally, in serverless runtims or long lived Kubernetes clusters.

## External Dependencies
None
