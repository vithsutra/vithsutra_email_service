//boiler plate code

package kafka

import (
	"context"
	"email-service/config"
	"email-service/internal/email"
	"email-service/internal/models"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer() {
	for _, topic := range config.Config.Kafka.Topics {
		go consumeTopic(topic)
	}
}

func consumeTopic(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Config.Kafka.Brokers,
		GroupID:  config.Config.Kafka.GroupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	log.Printf("Listening for messages on topic: %s\n", topic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		var emailData models.EmailRequest
		err = json.Unmarshal(msg.Value, &emailData)
		if err != nil {
			log.Println("Failed to unmarshal Kafka message:", err)
			continue
		}

		// Send Email
		err = email.SendEmail(emailData)
		if err != nil {
			log.Println("Failed to send email:", err)
		}
	}
}
