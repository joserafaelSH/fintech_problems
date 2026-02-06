package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/accounts"
	"github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/ledger"
	"github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/transactions"
	"github.com/joserafaelSH/fintech_problems/immutable_ledger_core/app/utils"
)

func main() {

	transactionDb := transactions.NewSafeTranasctionDatabase()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/transactions", transactions.CreateTransactionHandler(transactionDb, utils.GenerateID, utils.GenerateHash))
	r.GET("/accounts/:account_id/balances", accounts.GetAccountBalanceHanlder(transactionDb))
	r.GET("/ledger", ledger.GetLedger(transactionDb))
	r.GET("/ledger/verify", transactions.ValidateTransactionHandler(transactionDb))
	r.GET("/ledger/transactions", transactions.ListAllTransactions(transactionDb))

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
