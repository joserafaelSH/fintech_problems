package transactions

import (
	"fmt"
	"sync"
)

type TranasctionDatabase struct {
	store    map[string]TransactionModel
	lastHash string
	mut      sync.RWMutex
}

func NewSafeTranasctionDatabase() *TranasctionDatabase {
	return &TranasctionDatabase{
		store:    make(map[string]TransactionModel),
		lastHash: "echochain",
	}
}

func (db *TranasctionDatabase) Set(key string, value TransactionModel) {
	db.mut.Lock()
	defer db.mut.Unlock()
	db.store[key] = value
	db.lastHash = value.Hash
}

func (db *TranasctionDatabase) Get(key string) (TransactionModel, bool) {
	db.mut.RLock()
	defer db.mut.RUnlock()
	value, exists := db.store[key]
	return value, exists
}

func (db *TranasctionDatabase) GetDataFromAccount(accountId string) []TransactionModel {
	db.mut.RLock()
	defer db.mut.RUnlock()
	var transactions []TransactionModel
	for _, transaction := range db.store {
		if transaction.AccountId == accountId {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

func (db *TranasctionDatabase) GetAllTransactions(filters LedgerFilters) []TransactionModel {

	db.mut.RLock()
	defer db.mut.RUnlock()
	transactions := make([]TransactionModel, 0, len(db.store))
	for _, transaction := range db.store {
		if !match(transaction, filters) {
			continue
		}

		if filters.Limit != nil && len(transactions) >= *filters.Limit {
			break
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func match(tx TransactionModel, f LedgerFilters) bool {

	if f.AccountId != nil && tx.AccountId != *f.AccountId {
		return false
	}

	if f.AssetType != nil && tx.Asset.Unit != *f.AssetType {
		return false
	}

	if f.FromTimestamp != nil && tx.Timestamp.Before(*f.FromTimestamp) {
		fmt.Println("here")
		return false
	}

	if f.ToTimestamp != nil && tx.Timestamp.After(*f.ToTimestamp) {
		fmt.Println("here1")
		return false
	}

	return true
}
