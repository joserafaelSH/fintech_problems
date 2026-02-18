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

func SaveTransaction(tx Transaction) {
	insertSQL := `INSERT INTO transactions (id, account_id, amount, asset, created_at, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(insertSQL, tx.ID, tx.AccountID, tx.Amount, tx.Asset, tx.CreatedAt, tx.Status)
	if err != nil {
		logger.Logger.Error("Failed to save transaction", "error", err)
	}
}

func GetAllTransactions() ([]Transaction, error) {
	rows, err := database.DB.Query("SELECT id, account_id, amount, asset, created_at, status FROM transactions")
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
	Transaction Transaction
	Err         error
}

type TransactionProcessor struct {
	NumWorkers int
	InputChan  chan Transaction
	ResultChan chan Transaction
	ErrorChan  chan ErrorChannel
	GlobalCtx  context.Context
	Cancel     context.CancelFunc
	wg         sync.WaitGroup
	listenerWg sync.WaitGroup
}

func (t *Transaction) Process(r *rand.Rand) error {
	// simulate processing time
	time.Sleep(200 * time.Millisecond)

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

func CreateTransactionProcessor(parent context.Context) *TransactionProcessor {
	ctx, cancel := context.WithCancel(parent)
	return &TransactionProcessor{
		NumWorkers: NumWorkers,
		InputChan:  make(chan Transaction, InputChannelSize),
		ResultChan: make(chan Transaction, ResultChannelSize),
		ErrorChan:  make(chan ErrorChannel, ErrorChannelSize),
		GlobalCtx:  ctx,
		Cancel:     cancel,
	}
}

func (tp *TransactionProcessor) resultListener() {
	for tx := range tp.ResultChan {
		SaveTransaction(tx)
	}

}

func (tp *TransactionProcessor) errorListener() {
	for ec := range tp.ErrorChan {
		logger.Logger.Error("Error processing transaction", "transaction_id", ec.Transaction.ID, "error", ec.Err)
	}
}

func (tp *TransactionProcessor) worker(id int) {

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

			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
			err := tx.Process(r)
			if err != nil {
				select {
				case tp.ErrorChan <- ErrorChannel{Transaction: tx, Err: err}:
				default:
				}
			}

			select {
			case <-tp.GlobalCtx.Done():
				logger.Logger.Info("Worker shutting down due to global cancellation", "worker_id", id)
				return
			case tp.ResultChan <- tx:
				logger.Logger.Info("Worker processed transaction", "worker_id", id, "transaction_id", tx.ID)

			default:
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

	go func() {
		tp.resultListener()
	}()

	go func() {
		tp.errorListener()
	}()
}

func (tp *TransactionProcessor) Close() {
	tp.wg.Wait()
	tp.Cancel()
	close(tp.InputChan)
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
