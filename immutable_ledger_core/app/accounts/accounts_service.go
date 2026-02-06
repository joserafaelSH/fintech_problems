package accounts

import "github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/transactions"

func GetAccountBalanceService(transactionDb *transactions.TranasctionDatabase, accountId string) AccountBalance {

	accountData := transactionDb.GetDataFromAccount(accountId)
	balances := make(map[string]int64)
	for _, tx := range accountData {
		asset := tx.Asset
		if _, exists := balances[asset.Unit]; !exists {
			balances[asset.Unit] = asset.Amount
		} else {
			balances[asset.Unit] += asset.Amount
		}

	}

	return AccountBalance{
		AccountId: accountId,
		Balances:  balances,
	}
}
