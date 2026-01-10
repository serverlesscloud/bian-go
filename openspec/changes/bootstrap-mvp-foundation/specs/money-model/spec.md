# Money Model Capability

## ADDED Requirements

### Requirement: Decimal-Based Money Representation
The system SHALL represent financial amounts using decimal arithmetic (not floating-point) to prevent precision errors in financial calculations. Money amounts use the shopspring/decimal library for exact decimal representation.

#### Scenario: Money Creation
- **GIVEN** a financial amount of $123.45 AUD
- **WHEN** creating a Money instance
- **THEN** the amount is stored as decimal.Decimal ("123.45")
- **AND** the currency is stored as string "AUD" (ISO 4217 code)
- **AND** no precision loss occurs

#### Scenario: Floating-Point Prevention
- **GIVEN** the need to represent $0.10 + $0.20
- **WHEN** using Money type for calculation
- **THEN** the result is exactly $0.30
- **AND** NOT $0.30000000000000004 (floating-point error)
- **AND** all calculations maintain exact decimal precision

#### Scenario: Large Amount Handling
- **GIVEN** a very large amount like $999,999,999.99
- **WHEN** creating Money instance
- **THEN** the amount is stored exactly with no overflow
- **AND** arithmetic operations maintain precision
- **AND** no scientific notation in string representation

### Requirement: ISO 4217 Currency Codes
The system SHALL use ISO 4217 three-letter currency codes for all Money instances. Currency codes are validated against a defined set of supported currencies.

#### Scenario: Valid Currency Code
- **GIVEN** Money instance with currency "AUD"
- **WHEN** validating currency code
- **THEN** the currency is accepted (AUD is supported)
- **AND** no validation error occurs

#### Scenario: Supported Currencies
- **GIVEN** the Money type
- **WHEN** checking supported currencies
- **THEN** the system supports AUD (Australian Dollar)
- **AND** the system supports USD (US Dollar)
- **AND** the system supports GBP (British Pound)
- **AND** the system supports EUR (Euro)
- **AND** the system supports CAD (Canadian Dollar)

#### Scenario: Invalid Currency Code
- **GIVEN** attempting to create Money with currency "ABC"
- **WHEN** validating currency code
- **THEN** the system returns a validation error
- **AND** the error message indicates invalid currency code

#### Scenario: Currency Code Format
- **GIVEN** Money instance
- **WHEN** inspecting currency field
- **THEN** the currency code is uppercase three letters
- **AND** the currency code matches ISO 4217 standard
- **AND** lowercase currency codes are normalized to uppercase

### Requirement: Money Arithmetic Operations
The system SHALL provide arithmetic operations for Money type that maintain decimal precision and enforce currency compatibility. Operations include addition, subtraction, multiplication, and division.

#### Scenario: Addition - Same Currency
- **GIVEN** Money($100.00, AUD) and Money($50.00, AUD)
- **WHEN** adding the amounts
- **THEN** the result is Money($150.00, AUD)
- **AND** decimal precision is maintained
- **AND** no rounding errors occur

#### Scenario: Addition - Different Currency Error
- **GIVEN** Money($100.00, AUD) and Money($50.00, USD)
- **WHEN** attempting to add the amounts
- **THEN** the system returns an error
- **AND** the error indicates currency mismatch
- **AND** no result is returned (fail-fast)

#### Scenario: Subtraction - Same Currency
- **GIVEN** Money($100.00, AUD) and Money($30.00, AUD)
- **WHEN** subtracting the amounts
- **THEN** the result is Money($70.00, AUD)
- **AND** decimal precision is maintained
- **AND** the result can be negative if second amount is larger

#### Scenario: Multiplication by Scalar
- **GIVEN** Money($100.00, AUD) and scalar 2.5
- **WHEN** multiplying amount by scalar
- **THEN** the result is Money($250.00, AUD)
- **AND** decimal precision is maintained
- **AND** currency is preserved

#### Scenario: Division by Scalar
- **GIVEN** Money($100.00, AUD) and scalar 4
- **WHEN** dividing amount by scalar
- **THEN** the result is Money($25.00, AUD)
- **AND** decimal precision is maintained for exact division
- **AND** rounding mode is documented for non-exact division

#### Scenario: Division by Zero Error
- **GIVEN** Money($100.00, AUD)
- **WHEN** attempting to divide by zero
- **THEN** the system returns an error
- **AND** no result is returned

### Requirement: Money Comparison Operations
The system SHALL provide comparison operations for Money type including equality, greater than, less than, and zero checks. Comparisons enforce currency compatibility.

#### Scenario: Equality - Same Currency
- **GIVEN** Money($100.00, AUD) and Money($100.00, AUD)
- **WHEN** comparing for equality
- **THEN** the result is true (equal)

#### Scenario: Equality - Different Amount
- **GIVEN** Money($100.00, AUD) and Money($100.01, AUD)
- **WHEN** comparing for equality
- **THEN** the result is false (not equal)
- **AND** precision is exact (not approximate)

#### Scenario: Equality - Different Currency
- **GIVEN** Money($100.00, AUD) and Money($100.00, USD)
- **WHEN** comparing for equality
- **THEN** the result is false (different currencies)
- **OR** the system returns an error (currency mismatch)

#### Scenario: Greater Than Comparison
- **GIVEN** Money($100.00, AUD) and Money($50.00, AUD)
- **WHEN** comparing first > second
- **THEN** the result is true
- **AND** currency must match for comparison

#### Scenario: Less Than Comparison
- **GIVEN** Money($50.00, AUD) and Money($100.00, AUD)
- **WHEN** comparing first < second
- **THEN** the result is true

#### Scenario: Zero Check
- **GIVEN** Money($0.00, AUD)
- **WHEN** checking if amount is zero
- **THEN** the result is true
- **GIVEN** Money($0.01, AUD)
- **THEN** the result is false

#### Scenario: Negative Amount Check
- **GIVEN** Money($-100.00, AUD)
- **WHEN** checking if amount is negative
- **THEN** the result is true
- **AND** negative amounts are supported for representing debits

### Requirement: Money JSON Serialization
The system SHALL serialize Money type to JSON using string format for amounts to preserve decimal precision across JSON parsers. JSON deserialization reconstructs Money instances exactly.

#### Scenario: JSON Marshal
- **GIVEN** Money($123.45, AUD)
- **WHEN** marshaling to JSON
- **THEN** the output is {"amount": "123.45", "currency": "AUD"}
- **AND** amount is string format (not number)
- **AND** currency is uppercase string

#### Scenario: JSON Unmarshal
- **GIVEN** JSON {"amount": "123.45", "currency": "AUD"}
- **WHEN** unmarshaling to Money type
- **THEN** the result is Money($123.45, AUD) exactly
- **AND** decimal precision is preserved
- **AND** no floating-point conversion occurs

#### Scenario: JSON Unmarshal - Invalid Amount
- **GIVEN** JSON {"amount": "abc", "currency": "AUD"}
- **WHEN** unmarshaling to Money type
- **THEN** the system returns a parsing error
- **AND** the error message indicates invalid amount format

#### Scenario: JSON Unmarshal - Invalid Currency
- **GIVEN** JSON {"amount": "123.45", "currency": "XYZ"}
- **WHEN** unmarshaling to Money type
- **THEN** the system returns a validation error
- **AND** the error message indicates invalid currency code

#### Scenario: JSON Unmarshal - Number Format Prevention
- **GIVEN** JSON {"amount": 123.45, "currency": "AUD"} (amount as number, not string)
- **WHEN** unmarshaling to Money type
- **THEN** the system accepts and converts to string internally
- **OR** returns error requiring string format
- **AND** documentation warns against number format due to JSON parser precision limits

### Requirement: Money Model Structure
The system SHALL define Money as a Go struct with Amount (decimal.Decimal) and Currency (string) fields. The struct supports value semantics (copy by value).

#### Scenario: Money Struct Definition
- **GIVEN** the Money struct
- **WHEN** inspecting struct fields
- **THEN** the struct has Amount field of type decimal.Decimal
- **AND** the struct has Currency field of type string
- **AND** both fields are exported (capitalized)
- **AND** struct tags include JSON field names

#### Scenario: Money Constructor Function
- **GIVEN** the Money package
- **WHEN** creating new Money instance
- **THEN** a constructor function NewMoney(amount string, currency string) (Money, error) exists
- **AND** the constructor validates currency code
- **AND** the constructor parses amount string to decimal
- **AND** the constructor returns error for invalid inputs

#### Scenario: Money Zero Value
- **GIVEN** the Money type
- **WHEN** creating zero value Money{}
- **THEN** the amount is zero (decimal 0)
- **AND** the currency is empty string
- **AND** zero value is detectable (IsZero method)

### Requirement: Money Error Handling
The system SHALL provide clear error messages for all Money operations that can fail. Errors include currency mismatch, division by zero, invalid amounts, and invalid currency codes.

#### Scenario: Error Type Distinction
- **GIVEN** various Money operation failures
- **WHEN** examining error types
- **THEN** currency mismatch errors are identifiable
- **AND** validation errors are identifiable
- **AND** arithmetic errors are identifiable
- **AND** error messages are user-friendly

#### Scenario: Error Message Quality
- **GIVEN** an error from Money operation
- **WHEN** examining error message
- **THEN** the message clearly states what went wrong
- **AND** the message includes relevant details (e.g., "cannot add AUD and USD")
- **AND** the message does not expose internal implementation details
