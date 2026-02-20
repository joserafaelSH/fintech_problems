package rabbitmq

import (
	"context"

	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/database"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/transaction"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	chann     *amqp.Channel
	queue     amqp.Queue
	QueueName string
}

func CreateRabbitMQConnection() (*RabbitMQ, error) {

	rabbit := &RabbitMQ{QueueName: "transactions_queue"}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Logger.Error("Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}
	rabbit.conn = conn

	chann, err := rabbit.conn.Channel()
	if err != nil {
		logger.Logger.Error("Failed to open a channel", "error", err)
		return nil, err
	}
	rabbit.chann = chann
	// transactions_queue
	queue, err := chann.QueueDeclare(
		"transactions_queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		logger.Logger.Error("Failed to declare a queue", "error", err)
		return nil, err
	}
	rabbit.queue = queue
	logger.Logger.Info("RabbitMQ connection established and queue declared: ", "queue_name", queue.Name)

	return rabbit, nil
}

func (r *RabbitMQ) CloseRabbitMQConnection() {
	err := r.conn.Close()
	if err != nil {
		logger.Logger.Error("Failed to close RabbitMQ connection", "error", err)
	}

}

func (r *RabbitMQ) StartConsumer(database *database.Database) {

	msgs, err := r.chann.Consume(
		r.QueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		logger.Logger.Error("Failed to register a consumer", "error", err)
		panic(err)
	}

	var infinity chan struct{} = make(chan struct{})
	ctx := context.Background()
	tp := transaction.CreateTransactionProcessor(ctx, database)
	tp.Start()
	defer tp.Close()
	go func() {
		for data := range msgs {
			logger.Logger.Info("Received a message: ", "message", string(data.Body))
			tx, err := transaction.ParseTransaction(data.Body)
			if err != nil {
				logger.Logger.Error("Failed to parse transaction", "error", err)
				continue
			}
			tp.InputChan <- &tx
		}
	}()

	logger.Logger.Info("RabbitMQ consumer started, waiting for messages...")
	<-infinity

}
