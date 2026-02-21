package main

import (
	"context"

	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/database"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/rabbitmq"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/transaction"
)

const (
	workerPoolSize = 10
)

func main() {

	logger.CreateLogger()
	database, err := database.InitPostgresDb()
	if err != nil {
		logger.Logger.Error("Failed to initialize database", "error", err)
		return
	}
	defer database.DB.Close()
	logger.Logger.Info("Starting Transaction Worker Pool")

	ctx := context.Background()
	tp := transaction.NewProcessor(ctx, database, workerPoolSize)

	rmq, err := rabbitmq.CreateRabbitMQConnection()
	if err != nil {
		logger.Logger.Error("Failed to initialize RabbitMQ connection", "error", err)
		return
	}
	defer rmq.CloseRabbitMQConnection()
	rmq.StartConsumer(database, tp)

}
