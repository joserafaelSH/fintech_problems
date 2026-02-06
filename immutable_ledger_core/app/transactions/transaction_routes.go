package transactions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTransactionHandler(transactionDb *TranasctionDatabase, GenerateID func() string, GenerateHash func(string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionDto TransactionDto

		if err := c.BindJSON(&transactionDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := CreateTransaction(transactionDto, GenerateID, GenerateHash, transactionDb)
		if err != nil {
			var transErr TransactionError
			if errors.As(err, &transErr) {
				c.JSON(transErr.GetCode(), gin.H{
					"error": transErr.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Transaction created",
		})
	}

}

func ListAllTransactions(transactionDb *TranasctionDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		transactions := transactionDb.store
		c.JSON(http.StatusOK, gin.H{
			"transactions": transactions,
		})
	}

}

func ValidateTransactionHandler(transactionDb *TranasctionDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := ValidateTransactions(transactionDb)
		validity := hash != nil
		c.JSON(http.StatusOK, gin.H{
			"last_hash": hash,
			"valid":     validity,
		})
	}
}
