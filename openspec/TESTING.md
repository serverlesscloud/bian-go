
## 

#### Unit Tests

**Domain Layer:**
- Test interfaces with mock implementations
- Verify error handling
- Test edge cases (nil values, empty strings)

**Model Layer:**
- Test JSON serialization/deserialization
- Verify enum validation
- Test Money arithmetic operations

**GraphQL Layer:**
- Test resolvers with mock providers
- Verify query/mutation execution
- Test error mapping

#### Integration Tests

**Provider Tests:**
- Use `httptest` to mock external APIs
- Test full normalization pipeline
- Verify OAuth token refresh
- Test rate limiting behavior

**API Tests:**
- Execute actual REST requests against test server
- Execute actual GraphQL queries against test server
- Verify complete request/response cycle
- Test error responses
- Ensure API consistency (REST and GraphQL return equivalent data)

### Mock Provider

The `mock` provider serves dual purpose:
1. **Development** - Run server without external dependencies
2. **Testing** - Predictable test data for CI/CD

Includes realistic sample data:
- Multiple account types (checking, savings, credit card)
- Transaction history with various merchants
- Active, expired, and revoked consents
- Multi-currency balances
