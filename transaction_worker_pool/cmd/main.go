package main

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/database"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/rabbitmq"
)

func main() {

	logger.CreateLogger()
	database.InitDb()
	defer database.DB.Close()
	logger.Logger.Info("Starting Transaction Worker Pool")

	rmq := rabbitmq.CreateRabbitMQConnection()
	defer rmq.CloseRabbitMQConnection()
	rmq.StartConsumer()

}
