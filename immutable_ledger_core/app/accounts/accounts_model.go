package accounts

type AccountBalance struct {
	AccountId string           `json:"account_id"`
	Balances  map[string]int64 `json:"balances"`
}
