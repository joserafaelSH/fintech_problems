package transactions

type TransactionError interface {
	error
	GetCode() int
}

type TransactionConflictError struct {
	Message string
	Code    int
}

func (e *TransactionConflictError) Error() string {
	return e.Message
}

type TransactionNotFoundError struct {
	Message string
	Code    int
}

func (e *TransactionNotFoundError) Error() string {
	return e.Message
}

type TransactionValidationError struct {
	Message string
	Code    int
}

func (e *TransactionValidationError) Error() string {
	return e.Message
}

type TransactionRuleViolationError struct {
	Message string
	Code    int
}

func (e *TransactionRuleViolationError) Error() string {
	return e.Message
}

type TransactionMalformed struct {
	Message string
	Code    int
}

func (e *TransactionMalformed) Error() string {
	return e.Message
}
