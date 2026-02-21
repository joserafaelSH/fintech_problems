package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/transaction"
	_ "github.com/lib/pq"
)

const (
	insertTransactionRawSql  = `INSERT INTO transactions (id, account_id, amount, asset, created_at, status) VALUES (?, ?, ?, ?, ?, ?)`
	getAllTransactionsRawSql = `SELECT id, account_id, amount, asset, created_at, status FROM transactions`
)

type PostgreDatabase struct {
	DB                       *sql.DB
	InsertTransactionRawSql  string
	GetAllTransactionsRawSql string
}

func (db PostgreDatabase) createTransactionTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		account_id TEXT NOT NULL,
		amount BIGINT NOT NULL,
		asset TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		status TEXT NOT NULL
	);`

	_, err := db.DB.Exec(createTableSQL)
	if err != nil {
		logger.Logger.Error("Failed to create transactions table", "error", err)
	}
	logger.Logger.Info("Transactions table created or already exists")
}

func InitPostgresDb() (*PostgreDatabase, error) {

	dbConnection, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=transactions sslmode=disable")
	if err != nil {
		logger.Logger.Error("Failed to open database", "error", err)
		return nil, err
	}
	dbConnection.SetMaxOpenConns(25)
	dbConnection.SetMaxIdleConns(25)
	dbConnection.SetConnMaxLifetime(time.Hour)
	db := &PostgreDatabase{
		DB:                       dbConnection,
		InsertTransactionRawSql:  insertTransactionRawSql,
		GetAllTransactionsRawSql: getAllTransactionsRawSql,
	}

	db.createTransactionTable()
	return db, nil
}

func (p *PostgreDatabase) InsertTransaction(ctx context.Context, tx *transaction.Transaction) error {
	_, err := p.DB.ExecContext(ctx,
		`INSERT INTO transactions (id, account_id, amount, asset, created_at, status)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		tx.ID, tx.AccountID, tx.Amount, tx.Asset, tx.CreatedAt, tx.Status,
	)
	return err
}

func (p *PostgreDatabase) GetAllTransactions(ctx context.Context) ([]transaction.Transaction, error) {
	rows, err := p.DB.QueryContext(ctx,
		`SELECT id, account_id, amount, asset, created_at, status FROM transactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []transaction.Transaction

	for rows.Next() {
		var tx transaction.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Asset, &tx.CreatedAt, &tx.Status); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, rows.Err()
}
func (p *PostgreDatabase) Close() error {
	return p.DB.Close()
}
