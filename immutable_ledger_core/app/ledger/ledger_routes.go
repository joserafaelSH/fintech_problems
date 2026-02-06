package ledger

import (
	"github.com/gin-gonic/gin"
	"github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/transactions"
)

func GetLedger(transactionDb *transactions.TranasctionDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters transactions.LedgerFilters
		if err := c.BindQuery(&filters); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		transactionsList := transactionDb.GetAllTransactions(filters)
		c.JSON(200, gin.H{
			"transactions": transactionsList,
		})
	}
}
