package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/database"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
)

const (
	NumWorkers        = 100
	InputChannelSize  = 100
	ResultChannelSize = 100
	ErrorChannelSize  = 100
	PendingStatus     = "pending"
	CompletedStatus   = "completed"
	FailedStatus      = "failed"
)

type Transaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    int64     `json:"amount"`
	Asset     string    `json:"asset"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

func ParseTransaction(data []byte) (Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return tx, nil

}

type ErrorChannel struct {
	Transaction *Transaction
	Err         error
}

type TransactionProcessor struct {
	NumWorkers int
	InputChan  chan *Transaction
	ResultChan chan *Transaction
	ErrorChan  chan *ErrorChannel
	GlobalCtx  context.Context
	Cancel     context.CancelFunc
	wg         sync.WaitGroup
	listenerWg sync.WaitGroup
	Db         *database.Database
}

func (tp *TransactionProcessor) SaveTransaction(ctx context.Context, tx *Transaction) error {
	insertSQL := `INSERT INTO transactions (id, account_id, amount, asset, created_at, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := tp.Db.DB.ExecContext(ctx, insertSQL, tx.ID, tx.AccountID, tx.Amount, tx.Asset, tx.CreatedAt, tx.Status)
	if err != nil {
		logger.Logger.Error("Failed to save transaction", "error", err)
		return err
	}
	return nil
}

func (tp *TransactionProcessor) GetAllTransactions(ctx context.Context) ([]Transaction, error) {
	rows, err := tp.Db.DB.QueryContext(ctx, "SELECT id, account_id, amount, asset, created_at, status FROM transactions")
	if err != nil {
		logger.Logger.Error("Failed to query transactions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Asset, &tx.CreatedAt, &tx.Status)
		if err != nil {
			logger.Logger.Error("Failed to scan transaction row", "error", err)
			continue
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (t *Transaction) Process(ctx context.Context, r *rand.Rand) error {

	select {
	case <-time.After(200 * time.Millisecond):
	case <-ctx.Done():
		logger.Logger.Info("Transaction processing cancelled", "transaction_id", t.ID)
		return ctx.Err()
	}

	p := r.Float64()

	switch {
	case p < 0.2:
		t.Status = FailedStatus
		return fmt.Errorf("processing error for %s", t.ID)

	case p < 0.6:
		t.Status = FailedStatus
		return nil

	default:
		t.Status = CompletedStatus
		return nil
	}
}

func CreateTransactionProcessor(parent context.Context, db *database.Database) *TransactionProcessor {
	ctx, cancel := context.WithCancel(parent)
	return &TransactionProcessor{
		NumWorkers: NumWorkers,
		InputChan:  make(chan *Transaction, InputChannelSize),
		ResultChan: make(chan *Transaction, ResultChannelSize),
		ErrorChan:  make(chan *ErrorChannel, ErrorChannelSize),
		GlobalCtx:  ctx,
		Cancel:     cancel,
		Db:         db,
	}
}

func (tp *TransactionProcessor) resultListener() {
	defer tp.listenerWg.Done()
	for tx := range tp.ResultChan {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		tp.SaveTransaction(ctx, tx)
		cancel()
	}

}

func (tp *TransactionProcessor) errorListener() {
	defer tp.listenerWg.Done()
	for ec := range tp.ErrorChan {
		logger.Logger.Error("Error processing transaction", "transaction_id", ec.Transaction.ID, "error", ec.Err)
	}
}

func (tp *TransactionProcessor) worker(id int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
	defer tp.wg.Done()
	for {
		select {
		// Global context cancellation
		case <-tp.GlobalCtx.Done():
			logger.Logger.Info("Worker shutting down due to global cancellation", "worker_id", id)
			return

		case tx, ok := <-tp.InputChan:
			if !ok {
				logger.Logger.Warn("Worker shutting down due to input channel closure", "worker_id", id)
				return
			}

			ctx, cancel := context.WithTimeout(tp.GlobalCtx, 5*time.Second)
			err := tx.Process(ctx, r)
			cancel()
			if err != nil {

				tp.ErrorChan <- &ErrorChannel{Transaction: tx, Err: err}

			}

			select {
			case <-tp.GlobalCtx.Done():
				logger.Logger.Info("Worker shutting down due to global cancellation", "worker_id", id)
				return
			case tp.ResultChan <- tx:
				logger.Logger.Info("Worker processed transaction", "worker_id", id, "transaction_id", tx.ID)

			}

		}

	}
}

func (tp *TransactionProcessor) Start() {
	for i := 0; i < tp.NumWorkers; i++ {
		tp.wg.Add(1)
		go tp.worker(i)
	}
	tp.listenerWg.Add(2)
	go tp.resultListener()
	go tp.errorListener()
}

func (tp *TransactionProcessor) Close() {
	tp.Cancel()
	close(tp.InputChan)
	tp.wg.Wait()
	close(tp.ResultChan)
	close(tp.ErrorChan)
	tp.listenerWg.Wait()
}

func randomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CreateSampleTransactions() []Transaction {
	return []Transaction{
		{ID: randomString(), Amount: 100, Asset: "USD", CreatedAt: time.Now(), Status: PendingStatus},
		{ID: randomString(), Amount: 200, Asset: "EUR", CreatedAt: time.Now(), Status: PendingStatus},
		{ID: randomString(), Amount: 300, Asset: "JPY", CreatedAt: time.Now(), Status: PendingStatus},
	}
}
