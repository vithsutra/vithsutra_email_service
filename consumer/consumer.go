package consumer

import (
	"encoding/json"
	"log"

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
			if err := email.SendOtpEmail(emailRequest); err != nil {
				log.Println("error occurred while sending the otp email, Error: ", err.Error())
				continue
			}

			message.Ack(false)
		case "welcome":
			if err := email.WelcomeEmail(emailRequest); err != nil {
				log.Println("error occurred while sending the welocome email, Error: ", err.Error())
				continue
			}
			message.Ack(false)
		default:
			log.Println("invalid email type email request was receieved")
			message.Ack(false)
		}

	}
}
