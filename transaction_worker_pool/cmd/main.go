package main

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/database"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/rabbitmq"
)

func main() {

	logger.CreateLogger()
	database, err := database.InitDb()
	if err != nil {
		logger.Logger.Error("Failed to initialize database", "error", err)
		return
	}
	defer database.DB.Close()
	logger.Logger.Info("Starting Transaction Worker Pool")

	rmq, err := rabbitmq.CreateRabbitMQConnection()
	if err != nil {
		logger.Logger.Error("Failed to initialize RabbitMQ connection", "error", err)
		return
	}
	defer rmq.CloseRabbitMQConnection()
	rmq.StartConsumer(database)

}
