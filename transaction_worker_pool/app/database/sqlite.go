package database

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
)

const (
	InsertTransactionRawSql  = `INSERT INTO transactions (id, account_id, amount, asset, created_at, status) VALUES (?, ?, ?, ?, ?, ?)`
	GetAllTransactionsRawSql = `SELECT id, account_id, amount, asset, created_at, status FROM transactions`
)

type Database struct {
	DB                       *sql.DB
	InsertTransactionRawSql  string
	GetAllTransactionsRawSql string
}

func (db Database) createTransactionTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS transactions (
	"id" TEXT NOT NULL PRIMARY KEY,
	"account_id" TEXT,
	"amount" INTEGER,
	"asset" TEXT,
	"created_at" DATETIME,
	"status" TEXT
);`

	_, err := db.DB.Exec(createTableSQL)
	if err != nil {
		logger.Logger.Error("Failed to create transactions table", "error", err)
	}
	logger.Logger.Info("Transactions table created or already exists")
}

func InitDb() (*Database, error) {

	dbConnection, err := sql.Open("sqlite", "./transactions.db")
	if err != nil {
		logger.Logger.Error("Failed to open database", "error", err)
		return nil, err
	}

	db := &Database{
		DB:                       dbConnection,
		InsertTransactionRawSql:  InsertTransactionRawSql,
		GetAllTransactionsRawSql: GetAllTransactionsRawSql,
	}
	db.createTransactionTable()
	return db, nil
}
