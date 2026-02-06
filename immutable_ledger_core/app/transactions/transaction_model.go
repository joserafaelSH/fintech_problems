package transactions

import (
	"time"
)

type TransactionModel struct {
	TransactionId string    `json:"transaction_id"`
	AccountId     string    `json:"account_id"`
	Amount        int64     `json:"amount"`
	Timestamp     time.Time `json:"timestamp"`
	Asset         AssetType `json:"asset"`
	Hash          string    `json:"hash"`
	PreviousHash  string    `json:"previous_hash"`
}

type TransactionDto struct {
	AccountId string `json:"account_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"`
	Unit      string `json:"unit" binding:"required"`
}

type AssetType struct {
	Unit   string `json:"unit"`
	Amount int64  `json:"amount"`
}

func NewTransactionModel(transactionProperties TransactionDto, previousHash string, GenerateID func() string, GenerateHash func(string) string) TransactionModel {

	hashInput := transactionProperties.AccountId + previousHash
	return TransactionModel{
		TransactionId: GenerateID(),
		AccountId:     transactionProperties.AccountId,
		Amount:        transactionProperties.Amount,
		Asset:         AssetType{Unit: transactionProperties.Unit, Amount: transactionProperties.Amount},
		Timestamp:     time.Now().UTC(),
		Hash:          GenerateHash(hashInput),
		PreviousHash:  previousHash,
	}
}

type LedgerFilters struct {
	AccountId     *string    `form:"account_id" json:"account_id,omitempty" `
	AssetType     *string    `form:"asset_type" json:"asset_type,omitempty" `
	FromTimestamp *time.Time `form:"from_timestamp" json:"from_timestamp,omitempty" `
	ToTimestamp   *time.Time `form:"to_timestamp" json:"to_timestamp,omitempty" `
	Limit         *int       `form:"limit" json:"limit,omitempty" `
}
