package consumer

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/vithsutra/vithsutra_email_sending_service/email"
	"github.com/vithsutra/vithsutra_email_sending_service/internal/models"
	"github.com/vithsutra/vithsutra_email_sending_service/queue"
)

func Start() {
	conn, channel, queue := queue.Connect()

	defer conn.Close()
	defer channel.Close()

	messages, err := channel.Consume(
		queue.Name,
		"email-consumer-1",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("error occurred while connecing to queue, Error: ", err.Error())
	}

	log.Println("message service started succuessfully..")

	for message := range messages {
		emailRequest := new(models.Email)

		if err := json.Unmarshal(message.Body, emailRequest); err != nil {
			log.Println("error occurred while decoding the json messaage, Error: ", err.Error())
			message.Ack(false)
			continue
		}

		switch emailRequest.EmailType {
		case "otp":
			handleWithRetry(message, email.SendOtpEmail, emailRequest)
		case "welcome":
			handleWithRetry(message, email.WelcomeEmail, emailRequest)
		default:
			log.Println("Invalid email type received:", emailRequest.EmailType)
			message.Ack(false)
		}
	}
}

func handleWithRetry(message amqp.Delivery, sendFunc func(*models.Email) error, emailRequest *models.Email) {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		if err := sendFunc(emailRequest); err == nil {
			log.Printf("[%s] Email sent to: %s", emailRequest.EmailType, emailRequest.To)
			message.Ack(false)
			return
		} else {
			log.Printf("[%s] Retry %d/%d failed for: %s, Error: %v",
				emailRequest.EmailType, i+1, maxRetries, emailRequest.To, err)
		}
	}

	log.Printf("[%s] Email failed after %d retries: %s", emailRequest.EmailType, maxRetries, emailRequest.To)
	message.Nack(false, false)
}

