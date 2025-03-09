package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func Connect() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	url := os.Getenv("RABBITMQ_URL")

	if url == "" {
		log.Fatalln("missing RABBITMQ_URL env variable")
	}

	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatalln("error occurred while connecting to rabbitmq, Error: ", err.Error())
	}

	channel, err := conn.Channel()

	if err != nil {
		log.Fatalln("error occurred while opening the channel to rabbitmq, Error: ", err.Error())
	}

	queueName := os.Getenv("QUEUE_NAME")

	if queueName == "" {
		log.Fatalln("missing QUEUE_NAME env variable")
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-max-priority": int32(10),
		},
	)

	if err != nil {
		log.Fatalln("error occurred while creating the queue, Error: ", err.Error())
	}

	return conn, channel, &queue

}
