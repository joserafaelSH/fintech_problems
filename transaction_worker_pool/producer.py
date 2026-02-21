import pika
import json
import uuid
import random
import time
from datetime import datetime, timezone

RABBITMQ_HOST = "localhost"
QUEUE_NAME = "transactions_queue"

def random_transaction():
    return {
        "id": str(uuid.uuid4()),
        "account_id": f"acc_{random.randint(1, 1000)}",
        "amount": random.randint(-10000, 10000),  # cents, can simulate debit/credit
        "asset": random.choice(["USD", "EUR", "BTC", "ETH"]),
        "created_at": datetime.now(timezone.utc).isoformat(),
        "status": random.choice(["PENDING", "COMPLETED", "FAILED"]),
    }

def main(max_transactions=100):
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host=RABBITMQ_HOST)
    )
    channel = connection.channel()

    channel.queue_declare(queue=QUEUE_NAME, durable=False)

    print("Sending random transactions... Ctrl+C to stop.")

    try:
        count = 0
        while count < max_transactions:
            tx = random_transaction()
            message = json.dumps(tx)

            channel.basic_publish(
                exchange="",
                routing_key=QUEUE_NAME,
                body=message,
                properties=pika.BasicProperties(
                    delivery_mode=2  # make message persistent
                ),
            )

            print(f"Sent: {message}")
            count += 1
            #time.sleep(1)

    except KeyboardInterrupt:
        print("Stopping producer...")

    finally:
        connection.close()

if __name__ == "__main__":
    # read for command line argument for number of transactions to send
    import argparse
    parser = argparse.ArgumentParser(description="Transaction Producer")
    parser.add_argument("--max", type=int, default=100, help="Maximum number of transactions to send")
    args = parser.parse_args()
    print(f"Will send up to {args.max} transactions.")
    print("Press Ctrl+C to stop the producer.")
    if not args.max:
        print("No max transactions specified, running default (100).")
        args.max = 100
    main(max_transactions=args.max)
