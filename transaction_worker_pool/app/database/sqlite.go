package database

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
)

var DB *sql.DB

func InitDb() {

	db, err := sql.Open("sqlite", "./transactions.db")
	if err != nil {
		logger.Logger.Error("Failed to open database", "error", err)
	}

	DB = db
	CreateTransactionTable()
}

func CreateTransactionTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS transactions (
	"id" TEXT NOT NULL PRIMARY KEY,
	"account_id" TEXT,
	"amount" INTEGER,
	"asset" TEXT,
	"created_at" DATETIME,
	"status" TEXT
);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		logger.Logger.Error("Failed to create transactions table", "error", err)
	}
	logger.Logger.Info("Transactions table created or already exists")
}
