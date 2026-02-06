package transactions

import (
	"maps"
	"net/http"
	"slices"
	"sort"
)

func CreateTransaction(transactionDto TransactionDto, GenerateID func() string, GenerateHash func(string) string, transactionDb *TranasctionDatabase) error {

	if transactionDto.Amount == 0 {
		return &TransactionMalformed{
			Message: "Transaction amount cannot be zero",
			Code:    http.StatusUnprocessableEntity,
		}
	}

	previousHash := transactionDb.lastHash

	if previousHash == "echochain" && len(transactionDb.store) == 0 {
		previousHash = GenerateHash(previousHash + transactionDto.AccountId)
	}

	if previousHash == "echochain" && len(transactionDb.store) > 0 {
		return &TransactionRuleViolationError{
			Message: "Cannot create transaction: previous hash does not match the last transaction's hash",
			Code:    http.StatusInternalServerError,
		}
	}

	transactionModel := NewTransactionModel(transactionDto, previousHash, GenerateID, GenerateHash)
	if _, exists := transactionDb.Get(transactionModel.TransactionId); exists {
		return &TransactionConflictError{
			Message: "Transaction with the same ID already exists",
			Code:    http.StatusConflict,
		}
	}
	transactionDb.Set(transactionModel.TransactionId, transactionModel)
	return nil
}

func ValidateTransactions(transactionDb *TranasctionDatabase) *string {
	vals := slices.Collect(maps.Values(transactionDb.store))
	sort.Slice(vals, func(i, j int) bool {
		a := vals[i]
		b := vals[j]
		return a.Timestamp.Before(b.Timestamp)
	})

	for idx, transaction := range vals {
		if idx == 0 {
			continue
		}
		prevTransaction := vals[idx-1]
		if transaction.PreviousHash != prevTransaction.Hash {
			return nil
		}
	}
	lastHash := vals[len(vals)-1].Hash
	return &lastHash
}
