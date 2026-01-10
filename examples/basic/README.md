# Basic Example

This example demonstrates the BIAN-Go library with both REST and GraphQL APIs running on a single server.

## Running the Example

```bash
cd examples/basic
go run main.go
```

The server will start on port 8080 with the following endpoints:

- **REST API**: http://localhost:8080/
- **GraphQL API**: http://localhost:8080/graphql
- **GraphQL Playground**: http://localhost:8080/playground
- **Health Check**: http://localhost:8080/health

## REST API Examples

### Get Account Information
```bash
curl http://localhost:8080/accounts/acc-001
```

### Get Account Balance
```bash
curl http://localhost:8080/accounts/acc-001/balance
```

### Get All Account Balances
```bash
curl http://localhost:8080/accounts/acc-001/balances
```

### Get Account Transactions
```bash
# All transactions
curl http://localhost:8080/accounts/acc-001/transactions

# With date filtering
curl "http://localhost:8080/accounts/acc-001/transactions?fromDate=2024-01-01&toDate=2024-12-31"

# With pagination
curl "http://localhost:8080/accounts/acc-001/transactions?limit=5&offset=0"
```

### Get Transaction Details
```bash
curl http://localhost:8080/transactions/tx-001
```

### Get Consent Information
```bash
curl http://localhost:8080/consents/consent-001
```

### Get Consent Status
```bash
curl http://localhost:8080/consents/consent-001/status
```

## GraphQL Examples

Visit http://localhost:8080/playground to explore the GraphQL API interactively.

### Get Account Information
```graphql
query {
  account(id: "acc-001") {
    id
    accountNumber
    accountType
    productName
    nickname
    status
    openDate
    currency
  }
}
```

### Get Account Balance
```graphql
query {
  balance(accountId: "acc-001") {
    balanceType
    amount {
      amount
      currency
    }
    timestamp
  }
}
```

### Get All Account Balances
```graphql
query {
  balances(accountId: "acc-001") {
    balanceType
    amount {
      amount
      currency
    }
    timestamp
  }
}
```

### Get Account Transactions
```graphql
query {
  transactions(accountId: "acc-001", input: {
    limit: 10
    offset: 0
  }) {
    id
    transactionType
    amount {
      amount
      currency
    }
    description
    merchantName
    postingDate
    valueDate
  }
}
```

### Get Transaction Details
```graphql
query {
  transaction(id: "tx-001") {
    id
    reference
    transactionType
    amount {
      amount
      currency
    }
    description
    merchantName
    postingDate
    valueDate
    accountId
  }
}
```

### Get Consent Information
```graphql
query {
  consent(id: "consent-001") {
    id
    status
    scopes
    grantDate
    expiryDate
    revocationDate
  }
}
```

## Sample Data

The mock provider includes the following sample data:

### Accounts
- `acc-001`: Checking Account (AUD $2,547.83)
- `acc-002`: Savings Account (AUD $15,420.50)
- `acc-003`: Credit Card (USD -$1,250.75)

### Transactions
- Multiple transactions across all accounts
- Various transaction types (debit, credit, transfer, payment, fee)
- Realistic merchants and descriptions

### Consents
- `consent-001`: Active consent with full permissions
- `consent-002`: Expired consent
- `consent-003`: Revoked consent

## Configuration

The server can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `ENABLE_PLAYGROUND`: Enable GraphQL Playground (default: true)

Example:
```bash
PORT=3000 ENABLE_PLAYGROUND=false go run main.go
```