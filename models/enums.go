package models

// AccountType represents the type of bank account
type AccountType string

const (
	AccountTypeChecking   AccountType = "CHECKING"
	AccountTypeSavings    AccountType = "SAVINGS"
	AccountTypeCreditCard AccountType = "CREDIT_CARD"
	AccountTypeInvestment AccountType = "INVESTMENT"
)

// IsValid checks if the account type is valid
func (at AccountType) IsValid() bool {
	switch at {
	case AccountTypeChecking, AccountTypeSavings, AccountTypeCreditCard, AccountTypeInvestment:
		return true
	default:
		return false
	}
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDebit    TransactionType = "DEBIT"
	TransactionTypeCredit   TransactionType = "CREDIT"
	TransactionTypeTransfer TransactionType = "TRANSFER"
	TransactionTypePayment  TransactionType = "PAYMENT"
	TransactionTypeFee      TransactionType = "FEE"
)

// IsValid checks if the transaction type is valid
func (tt TransactionType) IsValid() bool {
	switch tt {
	case TransactionTypeDebit, TransactionTypeCredit, TransactionTypeTransfer, TransactionTypePayment, TransactionTypeFee:
		return true
	default:
		return false
	}
}

// ConsentStatus represents the status of a consent
type ConsentStatus string

const (
	ConsentStatusActive  ConsentStatus = "ACTIVE"
	ConsentStatusExpired ConsentStatus = "EXPIRED"
	ConsentStatusRevoked ConsentStatus = "REVOKED"
	ConsentStatusPending ConsentStatus = "PENDING"
)

// IsValid checks if the consent status is valid
func (cs ConsentStatus) IsValid() bool {
	switch cs {
	case ConsentStatusActive, ConsentStatusExpired, ConsentStatusRevoked, ConsentStatusPending:
		return true
	default:
		return false
	}
}

// BalanceType represents the type of balance
type BalanceType string

const (
	BalanceTypeCurrent   BalanceType = "CURRENT"
	BalanceTypeAvailable BalanceType = "AVAILABLE"
	BalanceTypePending   BalanceType = "PENDING"
)

// IsValid checks if the balance type is valid
func (bt BalanceType) IsValid() bool {
	switch bt {
	case BalanceTypeCurrent, BalanceTypeAvailable, BalanceTypePending:
		return true
	default:
		return false
	}
}

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusOpen      AccountStatus = "OPEN"
	AccountStatusClosed    AccountStatus = "CLOSED"
	AccountStatusSuspended AccountStatus = "SUSPENDED"
)

// IsValid checks if the account status is valid
func (as AccountStatus) IsValid() bool {
	switch as {
	case AccountStatusOpen, AccountStatusClosed, AccountStatusSuspended:
		return true
	default:
		return false
	}
}