package rabbitmq

import (
	"context"

	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/logger"
	"github.com/joserafaelSH/fintech_problems/transaction_worker_pool/app/transaction"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	chann     *amqp.Channel
	QueueName string
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Logger.Error("%s: %s", msg, err)
	}
}

func CreateRabbitMQConnection() *RabbitMQ {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	chann, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	queue, err := chann.QueueDeclare(
		"transactions_queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare a queue")
	logger.Logger.Info("RabbitMQ connection established and queue declared: ", "queue_name", queue.Name)

	return &RabbitMQ{conn: conn, chann: chann, QueueName: "transactions_queue"}
}

func (r *RabbitMQ) CloseRabbitMQConnection() {
	err := r.conn.Close()
	failOnError(err, "Failed to close RabbitMQ connection")
}

func (r *RabbitMQ) StartConsumer() {

	msgs, err := r.chann.Consume(
		r.QueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	var infinity chan struct{} = make(chan struct{})
	ctx := context.Background()
	tp := transaction.CreateTransactionProcessor(ctx)
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
			tp.InputChan <- tx
		}
	}()

	logger.Logger.Info("RabbitMQ consumer started, waiting for messages...")
	<-infinity
}
