package accounts

import (
	"github.com/gin-gonic/gin"
	"github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/transactions"
)

func GetAccountBalanceHanlder(transactionDb *transactions.TranasctionDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountId := c.Param("account_id")
		if accountId == "" {
			c.JSON(400, gin.H{
				"error": "account_id is required",
			})
			return
		}
		balance := GetAccountBalanceService(transactionDb, accountId)
		c.JSON(200, gin.H{
			"account_id": accountId,
			"balance":    balance,
		})
	}
}
